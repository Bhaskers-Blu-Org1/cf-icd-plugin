/**
 * Copyright IBM Corporation 2017.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

 package main

import (
    "fmt"
    "os"
    "bytes"
    "strings"
    "webhook"
    "os/exec"
    "encoding/json"
    "io/ioutil"
    "code.cloudfoundry.org/cli/plugin"
    "code.cloudfoundry.org/cli/plugin/models"
)

type ICDPlugin struct{}

type GitValues struct {
    GitURL string
    GitBranch string
    GitCommitID string
}

func GitInfo () ([]GitValues, error) {
    root_dir := ".git/refs/remotes/origin"
    files, err := ioutil.ReadDir(root_dir)
    var result []GitValues
    if err == nil {
        result = make([]GitValues, len(files))
        i := 0
        for _, file := range files {
            cmd := exec.Command("cat", root_dir + "/" + file.Name())
            var out1 bytes.Buffer
            cmd.Stdout = &out1
            err = cmd.Run()
            check(err)
            head := strings.Trim(out1.String(), "\n\r \b")
            cmd = exec.Command("git", "config", "--get", "remote.origin.url")
            var out2 bytes.Buffer
            cmd.Stdout = &out2
            err = cmd.Run()
            check(err)
            remote_url := strings.Trim(out2.String(), "\n\r \b")
            result[i] = GitValues {
                GitURL: remote_url,
                GitBranch: file.Name(),
                GitCommitID: head,
            }
            i += 1
        }
    }
    return result, nil
}

func (c *ICDPlugin) Run(cliConnection plugin.CliConnection, args []string) {
    var shouldRequest bool = false
    var method string = "POST"
    if args[0] == "icd" && len(args) == 4 && args[1] == "--create-connection" {
       shouldRequest = true
       method = "POST"
    } else if args[0] == "icd" && len(args) > 1 && args[1] == "--git-info" {
       res, err := GitInfo()
       check(err)
       fmt.Println(res)
    } else if args[0] == "icd" && len(args) == 4 && args[1] == "--delete-connection" {
       shouldRequest = true
       method = "DELETE"
    } else {
        fmt.Println("Enter a command. See cf icd --help for possible commands")
    }

    if (shouldRequest) {
        webhook_url := args[2]
        appName := args[3]
        fmt.Println("w: %s, a: %s", webhook_url, appName)
        current_org, err := cliConnection.GetCurrentOrg()
        check(err)
        current_space, err := cliConnection.GetCurrentSpace()
        check(err)
        apiEndpoint, err := cliConnection.ApiEndpoint()
        check(err)
        current_app, err := cliConnection.GetApp(appName)
        check(err)
        git_info, err := GitInfo()
        check(err)
        type Message struct {
            Org plugin_models.Organization
            Space plugin_models.Space
            App plugin_models.GetAppModel
            ApiEndpoint string
            Method string
            GitData []GitValues
        }
        amp := Message {
            Org: current_org,
            Space: current_space,
            App: current_app,
            ApiEndpoint: apiEndpoint,
            Method: method,
            GitData: git_info,
        }
        fmt.Println(amp)
        js, err := json.Marshal(amp)
        check(err)
        var buf = bytes.NewBufferString(string(js))

        body := webhook.Request(webhook_url, method, buf)
        fmt.Println(body)
    }
}

func (c *ICDPlugin) GetMetadata() plugin.PluginMetadata {
    return plugin.PluginMetadata{
        Name: "IBM Cloud Devops",
        Version: plugin.VersionType{
            Major: 0,
            Minor: 9,
            Build: 10,
        },
        MinCliVersion: plugin.VersionType{
            Major: 6,
            Minor: 7,
            Build: 0,
        },
        Commands: []plugin.Command{
            {
                Name:     "icd",
                HelpText: "IBM Cloud Devops plugin command's help text",

                // UsageDetails is optional
                // It is used to show help of usage of each command
                UsageDetails: plugin.Usage{
                    Usage: "cf icd [Options]",
					Options: map[string]string{
						"-create-connection [webhook_url] [app_name]": "Create a connection between a cf app and IBM toolchain",
						"-delete-connection [webhook_url] [app_name]": "Delete a connection",
					},
                },
            },
        },
    }
}

func check(e error) {
    if e != nil {
        fmt.Println("Warning: cf plugin icd operation not successful: %s",e)
        os.Exit(0)
        //panic(e)
    }
}

func main() {
    plugin.Start(new(ICDPlugin))
}
