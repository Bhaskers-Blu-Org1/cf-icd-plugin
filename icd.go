package main

import (
    "fmt"
    "bytes"
    "net/http"
    "encoding/json"
    "syscall"
    "os"
    "io/ioutil"
    "code.cloudfoundry.org/cli/plugin"
)

type ICDPlugin struct{}

func Request(url string, method string, buf *bytes.Buffer) (string) {
    client := &http.Client{}
    req, err := http.NewRequest(method, url, buf)
    check(err)
    resp, err := client.Do(req)
    check(err)
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    check(err)
    return string(body)
}

func WebhookConfigFile() (*os.File) {
    var webhookConfigFile = os.TempDir() + "webhook"
    var file *os.File
    var mode = os.FileMode(int(0600))
    if _, err := os.Stat(webhookConfigFile); os.IsNotExist(err) {
       file, err = os.Create(webhookConfigFile)
       check(err)
       err = (*file).Chmod(mode)
       check(err)
    } else {
       file, err = os.OpenFile(webhookConfigFile, syscall.O_RDWR, mode)
       check(err)
    }
    return file
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func WebhookConfig() (string) {
    var webhookConfigFile = os.TempDir() + "webhook"
    dat, err := ioutil.ReadFile(webhookConfigFile)
    check(err)
    return string(dat)
}

func (c *ICDPlugin) Run(cliConnection plugin.CliConnection, args []string) {
    if args[0] == "icd" && len(args) > 2 && args[1] == "--register-webhook" {
        var webhook = args[2]
        if webhook[:5] != "https" {
            fmt.Println("Error: https required");
            return;
        }
        var file = WebhookConfigFile()
        (*file).WriteString(webhook)
        err := (*file).Close()
        check(err)
    } else {
        output, err := cliConnection.CliCommand(args[1:]...);
        check(err)
        whcfg := WebhookConfig()
        fmt.Println(output)
        type Message struct {
            mtype string
            output []string
            args []string
        }
        amp := Message {
            mtype: "cf_command",
            output: output,
            args: args[1:],
        }
        js, err := json.Marshal(amp)
        check(err)
        fmt.Println(string(js))
        var buf = bytes.NewBufferString(string(js))

        Request(whcfg, "POST", buf)
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

func main() {
    plugin.Start(new(ICDPlugin))
}
