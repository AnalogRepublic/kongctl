package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/analogrepublic/kongctl/config"
	"github.com/analogrepublic/kongctl/data"
	"github.com/analogrepublic/kongctl/kong"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var api *kong.Kong

const (
	version     = "0.1.0"
	name        = "kongctl"
	description = "Kong management tool"
)

func main() {
	var err error

	config.Init()

	c := config.GetConfig()

	app := cli.NewApp()

	app.Name = name
	app.Usage = description
	app.Version = version
	app.EnableBashCompletion = true

	api, err = kong.NewKong(c.GetString("host"), c)

	if err != nil {
		fmt.Println(errors.Wrap(err, "Unable to communicate with the Kong service"))
		os.Exit(1)
	}

	app.Commands = []cli.Command{
		{
			Name:    "describe",
			Aliases: []string{"desc"},
			Usage:   "Describe a resource",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "enabled, e",
					Usage: "Only show enabled resources",
				},
			},
			Subcommands: []cli.Command{
				{
					Name:   "plugins",
					Usage:  "List plugins",
					Action: pluginListCommand,
				},
			},
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}

func pluginListCommand(c *cli.Context) error {
	params := &data.PluginRequestParams{}
	plugins, err := api.Plugins().List(params)

	if err != nil {
		return cli.NewExitError(errors.Wrap(err, "Unable to fetch the list of plugins"), 1)
	}

	table := tablewriter.NewWriter(os.Stdout)

	pluginsList := plugins.FilterData(func(p data.Plugin) bool {
		if c.GlobalBool("enabled") {
			return p.Enabled
		}

		return true
	})

	fmt.Printf("List of plugins, showing %d of %d \n", len(pluginsList), plugins.Total)
	table.SetHeader([]string{"ID", "API ID", "Name", "Enabled"})

	for _, plugin := range pluginsList {
		table.Append([]string{plugin.ID, plugin.ApiID, plugin.Name, fmt.Sprintf("%t", plugin.Enabled)})
	}

	table.Render()

	return nil
}
