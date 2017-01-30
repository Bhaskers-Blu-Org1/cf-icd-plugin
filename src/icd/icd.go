package main

import (
    "fmt"
    "bytes"
    "encoding/json"
    "icdlib"
    "code.cloudfoundry.org/cli/plugin"
)

type ICDPlugin struct{}

func (c *ICDPlugin) Run(cliConnection plugin.CliConnection, args []string) {
    if args[0] == "icd" && len(args) > 2 && args[1] == "--register-webhook" {
        var webhook = args[2]
        if webhook[:5] != "https" {
            fmt.Println("Error: https required");
            return;
        }
        var file = icdlib.WebhookConfigFile()
        (*file).WriteString(webhook)
        err := (*file).Close()
        icdlib.Check(err)
    } else {
        output, err := cliConnection.CliCommand(args[1:]...);
        icdlib.Check(err)
        whcfg := icdlib.WebhookConfig()
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
        icdlib.Check(err)
        fmt.Println(string(js))
        var buf = bytes.NewBufferString(string(js))

        icdlib.Request(whcfg, "POST", buf)
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
