package commands

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var Context = cli.Command{
	Name:    "context",
	Aliases: []string{"ctx", "domain"},
	Usage:   "Manage the context",
	Action:  showCommand,
	Subcommands: []cli.Command{
		{
			Name:   "switch",
			Usage:  "Switch to another context",
			Action: switchCommand,
		},
		{
			Name:   "list",
			Usage:  "List all contexts",
			Action: listCommand,
		},
	},
}

func showCommand(c *cli.Context) error {
	fmt.Printf("The current context is: %s\n", kongApi.Config.GetString("current_context"))
	return nil
}

func switchCommand(c *cli.Context) error {
	var contextName string

	// If we've provided an argument, assume
	// that the first one is the name.
	if c.NArg() > 0 {
		contextName = c.Args().Get(0)
	}

	// We need to pass in a contextName
	if contextName == "" {
		return cli.NewExitError(errors.New("You must provide a context to switch to"), 1)
	}

	if _, exists := kongApi.Config.FileData.Contexts[contextName]; !exists {
		return cli.NewExitError(errors.New("You must provide a valid context to switch to"), 1)
	}

	if contextName == kongApi.Config.FileData.CurrentContext {
		fmt.Printf("Not switching context, it is already %s\n", contextName)
		return nil
	}

	kongApi.Config.FileData.CurrentContext = contextName
	err := kongApi.Config.SaveFile()

	if err != nil {
		return cli.NewExitError(errors.New("Encountered an error when saving the config"), 1)
	}

	fmt.Printf("Context has been switched to: %s\n", contextName)

	return nil
}

func listCommand(c *cli.Context) error {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	table.SetCaption(true, "List of available contexts")
	table.SetHeader([]string{"Name", "Host", "Current"})

	for name, context := range kongApi.Config.FileData.Contexts {
		current := ""

		if name == kongApi.Config.FileData.CurrentContext {
			current = "yes"
		}

		table.Append([]string{name, context.Host, current})
	}

	table.Render()

	return nil
}
