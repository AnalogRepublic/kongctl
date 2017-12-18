package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/analogrepublic/kongctl/kong"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	kong, err := kong.NewKong("http://testing.kongctl.io:8001")

	if err != nil {
		fmt.Println(errors.Wrap(err, "Unable to communicate with the Kong service"))
		os.Exit(1)
	}

	app.Name = "kongctl"
	app.Usage = "Kong management tool"
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		{
			Name:    "plugins",
			Aliases: []string{"p"},
			Usage:   "Manage plugins",
			Action: func(c *cli.Context) error {
				plugins, err := kong.Plugins().List()

				if err != nil {
					return cli.NewExitError(errors.Wrap(err, "Unable to fetch the list of plugins"), 1)
				}

				table := tablewriter.NewWriter(os.Stdout)

				fmt.Printf("List of plugins, showing %d of %d \n", len(plugins.Data), plugins.Total)
				table.SetHeader([]string{"ID", "API ID", "Name", "Enabled"})

				for _, plugin := range plugins.Data {
					table.Append([]string{plugin.ID, plugin.ApiID, plugin.Name, fmt.Sprintf("%t", plugin.Enabled)})
				}

				table.Render()

				return nil
			},
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}
