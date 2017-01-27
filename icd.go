package main

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "code.cloudfoundry.org/cli/plugin"
)

type ICDPlugin struct{}

func Request(url string, cliConnection plugin.CliConnection) (*map[string]interface{}, error) {
    at, aterr := cliConnection.AccessToken();
    var dat map[string]interface{}
    if aterr != nil {
       fmt.Println("AT Err: ", aterr);
       return &dat, aterr;
    } else {
       client := &http.Client{}
       req, err := http.NewRequest("GET", url, nil)
       if err != nil {
          return &dat, err
       }
       req.Header.Add("Authorization", at)
       resp, err := client.Do(req)
       if err != nil {
          return &dat, err
       }
       defer resp.Body.Close()
       body, err := ioutil.ReadAll(resp.Body)
       if err != nil {
          return &dat, err
       }
       errs := json.Unmarshal(body, &dat);
       if errs != nil {
          return &dat, errs
       }
    }
    return &dat, nil
}

func (c *ICDPlugin) Run(cliConnection plugin.CliConnection, args []string) {
    var appname = ""
    var tcid = ""
    var args2 = make([]string,0)
    for idx := range args {
       if args[idx] == "--create-deployable-mapping" {
          tcid = args[idx+1]
          idx += 1
       } else {
          args2 = append(args2[:len(args2)], args[:idx-1]...)
       }
    }
    fmt.Println(args2)

    output, err := cliConnection.CliCommand(args2[1:]...);
    if err != nil {
      fmt.Println("PLUGIN OUTPUT: Output from CliCommand: ", output)
      fmt.Println("PLUGIN ERROR: Error from CliCommand: ", err)
    }

    if args[0] == "icd" && len(args2) > 1 && tcid != "" {
       output, err := cliConnection.CliCommand("app " + appname + "--guid");
       fmt.Println(output)
       fmt.Println(err)
       urlbase := "https://otc-api.stage1.ng.bluemix.net/api/v1/";
       dat, err := Request(urlbase + "toolchains/" + tcid + "/services", cliConnection)
       if err != nil {
          fmt.Println("err: ", err)
       }
       fmt.Println(*dat)
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
