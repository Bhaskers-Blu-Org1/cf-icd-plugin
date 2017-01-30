package main

import (
    "fmt"
    "bytes"
    "encoding/json"
    "webhook"
    "code.cloudfoundry.org/cli/plugin"
)

type ICDPlugin struct{}

func (c *ICDPlugin) Run(cliConnection plugin.CliConnection, args []string) {
    if args[0] == "icd" && len(args) > 2 && args[1] == "--register-webhook" {
        var webhook_url = args[2]
        if webhook_url[:5] != "https" {
            fmt.Println("Error: https required");
            return;
        }
        var file = webhook.ConfigFile()
        (*file).WriteString(webhook_url)
        err := (*file).Close()
        check(err)
    } else {
        output, err := cliConnection.CliCommand(args[1:]...);
        check(err)
        whcfg := webhook.Config()
        fmt.Println(output)
        type Message struct {
            Mtype string
            Output []string
            Args []string
        }
        amp := Message {
            Mtype: "cf_command",
            Output: output,
            Args: args[1:],
        }
        fmt.Println(amp)
        js, err := json.Marshal(amp)
        check(err)
        fmt.Println(js)
        var buf = bytes.NewBufferString(string(js))

        webhook.Request(whcfg, "POST", buf)
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
