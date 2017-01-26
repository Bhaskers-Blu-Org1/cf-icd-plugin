package main

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "code.cloudfoundry.org/cli/plugin"
)

// BasicPlugin is the struct implementing the interface defined by the core CLI. It can
// be found at  "code.cloudfoundry.org/cli/plugin/plugin.go"
type ICDPlugin struct{}

// Run must be implemented by any plugin because it is part of the
// plugin interface defined by the core CLI.
//
// Run(....) is the entry point when the core CLI is invoking a command defined
// by a plugin. The first parameter, plugin.CliConnection, is a struct that can
// be used to invoke cli commands. The second paramter, args, is a slice of
// strings. args[0] will be the name of the command, and will be followed by
// any additional arguments a cli user typed in.
//
// Any error handling should be handled with the plugin itself (this means printing
// user facing errors). The CLI will exit 0 if the plugin exits 0 and will exit
// 1 should the plugin exits nonzero.
func (c *ICDPlugin) Run(cliConnection plugin.CliConnection, args []string) {
    // Ensure that we called the command basic-plugin-command
    if args[0] == "icd" && len(args) > 1{
       fmt.Println("Running the IBM Continuous Delivery command");
           at, aterr := cliConnection.AccessToken();
           if aterr != nil {
              fmt.Println("AT Err: ", aterr);
           } else {
              client := &http.Client{}

              url := "https://otc-api.stage1.ng.bluemix.net/api/v1/toolchains";
              //resp, err := client.Get(url);

              req, err := http.NewRequest("GET", url, nil)
              req.Header.Add("Authorization", at)
              resp, err := client.Do(req)
              defer resp.Body.Close()
              body, err := ioutil.ReadAll(resp.Body)
              var dat map[string]interface{}
              errs := json.Unmarshal(body, &dat);
 
              fmt.Println(dat["total_results"])
              strs := dat["items"].([]interface{})
              str1 := strs[0].(map[string]interface {})
              fmt.Println(str1["toolchain_guid"])
              //for idx, val := range dat["items"] {
              //    fmt.Println(val["toolchain_guid"])
              //}
              fmt.Println(errs)
              fmt.Println(err)
           }
           output, err := cliConnection.CliCommand(args[1:]...);
           // The call to plugin.CliCommand() returns an error if the cli command
       // returns a non-zero return code. The output written by the CLI
       // is returned in any case.
       if err != nil {
        fmt.Println("PLUGIN OUTPUT: Output from CliCommand: ", output)
        fmt.Println("PLUGIN ERROR: Error from CliCommand: ", err)
       }

    }
}

// GetMetadata must be implemented as part of the plugin interface
// defined by the core CLI.
//
// GetMetadata() returns a PluginMetadata struct. The first field, Name,
// determines the name of the plugin which should generally be without spaces.
// If there are spaces in the name a user will need to properly quote the name
// during uninstall otherwise the name will be treated as seperate arguments.
// The second value is a slice of Command structs. Our slice only contains one
// Command Struct, but could contain any number of them. The first field Name
// defines the command `cf basic-plugin-command` once installed into the CLI. The
// second field, HelpText, is used by the core CLI to display help information
// to the user in the core commands `cf help`, `cf`, or `cf -h`.
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
