package commands

import (
	"fmt"
	"os"

	"github.com/analogrepublic/kongctl/data"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var Describe = cli.Command{
	Name:    "describe",
	Aliases: []string{"desc", "inspect"},
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
			Action: describePluginCommand,
		},
		{
			Name:   "apis",
			Usage:  "List apis",
			Action: describeApiCommand,
		},
	},
}

func describeApiCommand(c *cli.Context) error {
	params := &data.ApiRequestParams{}
	table := tablewriter.NewWriter(os.Stdout)

	// Style the table
	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
	table.SetCenterSeparator("|")

	// By default we'll not specify a name to find.
	name := ""

	// If we've provided an argument, assume
	// that the first one is the name.
	if c.NArg() > 0 {
		name = c.Args().Get(0)
	}

	// If we've not passed in a name as an argument,
	// we will just fetch a list of all apis.
	if name == "" {
		apiList, err := kongApi.Apis().List(params)

		if err != nil {
			return cli.NewExitError(errors.Wrap(err, "Unable to fetch the list of apis"), 1)
		}

		table.SetCaption(true, fmt.Sprintf("List of apis, showing %d of %d \n", len(apiList.Data), apiList.Total))
		table.SetHeader([]string{"ID", "Name", "Upstream URL"})

		for _, api := range apiList.Data {
			table.Append([]string{api.ID, api.Name, api.UpstreamUrl})
		}
	} else {
		params.ID = name
		apiItem, err := kongApi.Apis().Retrieve(params)

		if err != nil {
			return cli.NewExitError(errors.Wrap(err, fmt.Sprintf("Unable to fetch the api with the ID %s", name)), 1)
		}

		if apiItem.ID == "" {
			return cli.NewExitError(errors.New("Could not find api with the requested ID"), 1)
		}

		table.SetHeader([]string{"Key", "Value"})

		data := [][]string{
			{"ID", apiItem.ID},
			{"Name", apiItem.Name},
		}

		table.SetAutoMergeCells(true)
		table.SetRowLine(true)
		table.AppendBulk(data)
	}

	table.Render()

	return nil
}

func describePluginCommand(c *cli.Context) error {
	params := &data.PluginRequestParams{}
	plugins, err := kongApi.Plugins().List(params)

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

	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	table.SetCaption(true, fmt.Sprintf("List of plugins, showing %d of %d \n", len(pluginsList), plugins.Total))
	table.SetHeader([]string{"ID", "API ID", "Name", "Enabled"})

	for _, plugin := range pluginsList {
		table.Append([]string{plugin.ID, plugin.ApiID, plugin.Name, fmt.Sprintf("%t", plugin.Enabled)})
	}

	table.Render()

	return nil
}
