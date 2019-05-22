package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"os"
)

type RenamerPlugin struct{}

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
func (c *RenamerPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "renamify" {
		appName := args[1]
		newName := appName + "-potato"

		fmt.Printf("Renaming app '%s' to '%s'\n\n", appName, newName)

		app, err := cliConnection.GetApp(appName)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		path := fmt.Sprintf("/v3/apps/%s", app.Guid)
		requestBody := fmt.Sprintf(`{"name": "%s"}`, newName)

		_, err = cliConnection.CliCommandWithoutTerminalOutput("curl", path, "-X", "PATCH", "-d", requestBody)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		fmt.Printf("Renamed your app!\n\n")

		_, err = cliConnection.CliCommand("app", newName)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
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
func (c *RenamerPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "RenamerPlugin",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "renamify",
				HelpText: "Append '-potato' to an app name, because potatoes are awesome!",
				UsageDetails: plugin.Usage{
					Usage: "cf renamify APP_NAME",
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
	plugin.Start(new(RenamerPlugin))
	// Plugin code should be written in the Run([]string) method,
	// ensuring the plugin environment is bootstrapped.
}