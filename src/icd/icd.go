package main

import (
    "fmt"
    "bytes"
    "encoding/json"
    "webhook"
    "code.cloudfoundry.org/cli/plugin"
    "code.cloudfoundry.org/cli/plugin/models"
)

type ICDPlugin struct{}

func (c *ICDPlugin) Run(cliConnection plugin.CliConnection, args []string) {
    if args[0] == "icd" && len(args) > 2 && args[1] == "--create-connection" {
        tcName := args[2]
        appName := args[3]
        fmt.Println("t: %s, a: %s", tcName, appName)
        whcfg, err := webhook.Config()
        current_org, err := cliConnection.GetCurrentOrg()
        check(err)
        current_space, err := cliConnection.GetCurrentSpace()
        check(err)
        apiEndpoint, err := cliConnection.ApiEndpoint()
        check(err)
        current_apps, err := cliConnection.GetApps()
        check(err)
        current_token, err := cliConnection.AccessToken()
        check(err)
        type Message struct {
            Org plugin_models.Organization
            Space plugin_models.Space
            Apps []plugin_models.GetAppsModel
            ApiEndpoint string
            AccessToken string
        }
        amp := Message {
            Org: current_org,
            Space: current_space,
            Apps: current_apps,
            ApiEndpoint: apiEndpoint,
            AccessToken: current_token,
        }
        fmt.Println(amp.Org)
        js, err := json.Marshal(amp)
        check(err)
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
