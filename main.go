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

var kongApi *kong.Kong

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

	kongApi, err = kong.NewKong(c.GetString("host"), c)

	if err != nil {
		fmt.Println(errors.Wrap(err, "Unable to communicate with the Kong service"))
		os.Exit(1)
	}

	app.Commands = []cli.Command{
		{
			Name:    "apply",
			Aliases: []string{"upload", "set", "sync"},
			Usage:   "Apply a yaml file to a Kong service",
			Action: func(c *cli.Context) error {
				fmt.Println("Should be applyin' shit here")

				// Specify the file we're reading
				yamlFile := &data.UnparsedYamlFile{
					Name: "test_files/file.yaml",
				}

				// Convert it to a struct
				parsedFile, err := yamlFile.Unmarshal()

				if err != nil {
					return cli.NewExitError(err, 1)
				}

				// Grab the first Api from that struct
				testApiItem := parsedFile.Apis[0]

				// Create it
				createdApiItem, err := kongApi.Apis().Add(testApiItem)

				if err != nil {
					return cli.NewExitError(err, 1)
				}

				if createdApiItem.ID == "" {
					// We didn't get a new one, so it must exist, let's update it instead
					updatedApiItem, err := kongApi.Apis().Update(&data.ApiRequestParams{Name: testApiItem.Name}, testApiItem)

					testApiItem = updatedApiItem

					if err != nil {
						return cli.NewExitError(err, 1)
					}
				} else {
					testApiItem = createdApiItem
				}

				// Create a fake parsed file with the api we've just got back
				fakeParsedFile := &data.ParsedYamlFile{
					Apis: []*data.Api{
						testApiItem,
					},
					Plugins: []*data.Plugin{},
				}

				// Convert our fake parsed object to a string
				fakeYamlOutput, err := fakeParsedFile.Marshal()

				if err != nil {
					return cli.NewExitError(err, 1)
				}

				fmt.Println(fakeYamlOutput)

				return nil
			},
		},
		{
			Name:    "export",
			Aliases: []string{"download", "fetch", "grab"},
			Usage:   "Export resources to a yaml file",
			Action: func(c *cli.Context) error {
				fmt.Println("Should be exportin' shit here")

				return nil
			},
		},
		{
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
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
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
