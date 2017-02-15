package main

import (
    "fmt"
    "log"
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

func GitInfo () (GitValues, error) {
    _, err := ioutil.ReadDir(".git")
    var result GitValues
    if err == nil {
        head, err := ioutil.ReadFile(".git/HEAD")
        check(err)
        fmt.Println(string(head))
        parts := strings.Split(string(head), "/")
        branch_name := strings.Trim(parts[len(parts) - 1], "\n\r \b")
        fmt.Println(branch_name)
        id, err := ioutil.ReadFile(".git/refs/heads/" + branch_name)
        check(err)
        fmt.Println(string(id))
        cmd := exec.Command("git", "config", "--get", "remote.origin.url")
        var out bytes.Buffer
        cmd.Stdout = &out
        err = cmd.Run()
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("in all caps: %q\n", out.String())
        remote_url := strings.Trim(out.String(), "\n\r \b")
        commit_id := strings.Trim(string(id), "\n\r \b")
        result = GitValues {
            GitURL: remote_url,
            GitBranch: branch_name,
            GitCommitID: commit_id,
        }
    }
    return result, nil
}

func (c *ICDPlugin) Run(cliConnection plugin.CliConnection, args []string) {
    var shouldRequest bool = false
    var method string = "POST"
    if args[0] == "icd" && len(args) > 2 && args[1] == "--create-connection" {
       shouldRequest = true
       method = "POST"
    } else if args[0] == "icd" && len(args) > 1 && args[1] == "--git-info" {
       res, err := GitInfo()
       check(err)
       fmt.Println(res)
    } else if args[0] == "icd" && len(args) > 2 && args[1] == "--delete-connection" {
       shouldRequest = true
       method = "DELETE"
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
            GitData GitValues
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
        Name: "IBM Continuous Delivery",
        Version: plugin.VersionType{
            Major: 0,
            Minor: 0,
            Build: 1,
        },
        MinCliVersion: plugin.VersionType{
            Major: 6,
            Minor: 7,
            Build: 0,
        },
        Commands: []plugin.Command{
            {
                Name:     "icd",
                HelpText: "IBM Continous Delivery plugin command's help text",

                // UsageDetails is optional
                // It is used to show help of usage of each command
                UsageDetails: plugin.Usage{
                    Usage: "IBM Continous Delivery:\n   cf icd",
                },
            },
        },
    }
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    plugin.Start(new(ICDPlugin))
}
