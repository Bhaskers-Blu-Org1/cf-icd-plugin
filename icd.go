package main

import (
    "fmt"
    "net/http"
    "syscall"
    "os"
    "io/ioutil"
    "code.cloudfoundry.org/cli/plugin"
)

type ICDPlugin struct{}

func Request(url string) (*map[string]interface{}, error) {
    var dat map[string]interface{}
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
       return &dat, err
    }
    resp, err := client.Do(req)
    if err != nil {
       return &dat, err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
       return &dat, err
    }
    fmt.Println(body)
    return &dat, nil
}

func WebhookConfigFile() (*os.File, error) {
    var webhookConfigFile = os.TempDir() + "webhook"
    var file *os.File
    var mode = os.FileMode(int(0600))
    var err error
    if _, err := os.Stat(webhookConfigFile); os.IsNotExist(err) {
       file, err = os.Create(webhookConfigFile)
       if err != nil {
          return nil, err
       }
       err = (*file).Chmod(mode)
    } else {
       file, err = os.OpenFile(webhookConfigFile, syscall.O_RDWR, mode)
    }
    return file, err
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
    fmt.Print(string(dat))
    return string(dat)
}

func (c *ICDPlugin) Run(cliConnection plugin.CliConnection, args []string) {
    if args[0] == "icd" && len(args) > 2 && args[1] == "--register-webhook" {
        var webhook = args[2]
        if webhook[:5] != "https" {
            fmt.Println("Error: https required");
            return;
        }
        var file, err = WebhookConfigFile()
        if err != nil {
           fmt.Println("Error: ", err)
           return
        }
        (*file).WriteString(webhook)
        err = (*file).Close()
        fmt.Println("Error: ", err)
        if err != nil {
           fmt.Println("Error: ", err)
           return
        }
    } else {
        output, err := cliConnection.CliCommand(args[1:]...);
        if err != nil {
          fmt.Println("PLUGIN OUTPUT: Output from CliCommand: ", output)
          fmt.Println("PLUGIN ERROR: Error from CliCommand: ", err)
        }
        whcfg := WebhookConfig()
        if err != nil {
           fmt.Println("Error: ", err)
           return
        }
        fmt.Println(whcfg)
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

// Unlike most Go programs, the `Main()` function will not be used to run all of the
// commands provided in your plugin. Main will be used to initialize the plugin
// process, as well as any dependencies you might require for your
// plugin.
func main() {
    // Any initialization for your plugin can be handled here
    //
    // Note: to run the plugin.Start method, we pass in a pointer to the struct
    // implementing the interface defined at "code.cloudfoundry.org/cli/plugin/plugin.go"
    //
    // Note: The plugin's main() method is invoked at install time to collect
    // metadata. The plugin will exit 0 and the Run([]string) method will not be
    // invoked.
    plugin.Start(new(ICDPlugin))
    // Plugin code should be written in the Run([]string) method,
    // ensuring the plugin environment is bootstrapped.
}
