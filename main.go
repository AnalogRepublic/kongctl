package kongctl

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/analogrepublic/kongctl/commands"
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

	err = config.Init()

	if err != nil {
		fmt.Println(errors.Wrap(err, "Unable to read config file"))
		os.Exit(1)
	}

	c := config.GetConfig()

	app := cli.NewApp()

	app.Name = name
	app.Usage = description
	app.Version = version
	app.EnableBashCompletion = true

	context, err := c.GetCurrentContext()

	if err != nil {
		fmt.Println(errors.Wrap(err, "Unable to determine context to use, please check your config file"))
		os.Exit(1)
	}

	kongApi, err = kong.NewKong(context.Host, c)

	if err != nil {
		fmt.Println(errors.Wrap(err, "Unable to communicate with the Kong service"))
		os.Exit(1)
	}

	commands.SetKongApi(kongApi)

	app.Commands = []cli.Command{
		{
			Name:    "context",
			Aliases: []string{"ctx", "domain"},
			Usage:   "Manage the context",
			Action:  contextShowCommand,
			Subcommands: []cli.Command{
				{
					Name:   "switch",
					Usage:  "Switch to another context",
					Action: contextSwitchCommand,
				},
				{
					Name:   "list",
					Usage:  "List all contexts",
					Action: contextListCommand,
				},
			},
		},
		{
			Name:    "apply",
			Aliases: []string{"upload", "set", "sync"},
			Usage:   "Apply a yaml file to a Kong service",
			Action:  applyCommand,
		},
		{
			Name:    "export",
			Aliases: []string{"download", "fetch", "grab"},
			Usage:   "Export resources to a yaml file",
			Action:  exportCommand,
		},
		commands.Describe,
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}

func applyCommand(c *cli.Context) error {
	var file string

	// If we've provided an argument, assume
	// that the first one is the name.
	if c.NArg() > 0 {
		file = c.Args().Get(0)
	}

	// We need to pass in a file
	if file == "" {
		return cli.NewExitError(errors.New("You must provide a file to apply"), 1)
	}

	// Specify the file we're reading
	yamlFile := &data.UnparsedYamlFile{
		Name: file,
	}

	// Convert it to a struct
	parsedFile, err := yamlFile.Unmarshal()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	remoteApis, err := kongApi.Apis().List(&data.ApiRequestParams{})

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	remotePlugins, err := kongApi.Plugins().List(&data.PluginRequestParams{})

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	// Create a parsed file with the apis we've just got back from
	// the remote kong service
	remoteParsedFile := &data.ParsedYamlFile{
		Apis:    remoteApis.Data,
		Plugins: remotePlugins.Data,
	}

	diff, err := parsedFile.Diff(remoteParsedFile)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	for _, addition := range diff.Additions {
		apiAddition := addition.(data.Api)

		fmt.Printf("Adding api '%s'\n", apiAddition.Name)

		result, err := kongApi.Apis().Add(&apiAddition)

		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if result.ID == "" {
			return cli.NewExitError(errors.New(fmt.Sprintf("Failed creating api endpoint %s", apiAddition.Name)), 1)
		}
	}

	for _, updatePair := range diff.Updates {
		update := updatePair.([]interface{})[0]
		apiUpdate := update.(data.Api)

		fmt.Printf("Updating api '%s'\n", apiUpdate.Name)

		result, err := kongApi.Apis().Update(&data.ApiRequestParams{Name: apiUpdate.Name}, &apiUpdate)

		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if result.ID == "" {
			return cli.NewExitError(errors.New(fmt.Sprintf("Failed updating api endpoint %s", apiUpdate.Name)), 1)
		}
	}

	for _, deletion := range diff.Deletions {
		apiDeletion := deletion.(data.Api)
		fmt.Printf("Removing api '%s'\n", apiDeletion.Name)
		err := kongApi.Apis().Delete(&data.ApiRequestParams{Name: apiDeletion.Name})

		if err != nil {
			return cli.NewExitError(errors.New(fmt.Sprintf("Failed removing api endpoint %s", apiDeletion.Name)), 1)
		}
	}

	return nil
}

func exportCommand(c *cli.Context) error {
	var file string

	// If we've provided an argument, assume
	// that the first one is the name.
	if c.NArg() > 0 {
		file = c.Args().Get(0)
	}

	// We need to pass in a file
	if file == "" {
		return cli.NewExitError(errors.New("You must provide a file to export the kong services to"), 1)
	}

	remoteApis, err := kongApi.Apis().List(&data.ApiRequestParams{})

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	remotePlugins, err := kongApi.Plugins().List(&data.PluginRequestParams{})

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	remoteParsedFile := &data.ParsedYamlFile{
		Apis:    remoteApis.Data,
		Plugins: remotePlugins.Data,
	}

	remoteParsedFileOutput, err := remoteParsedFile.Marshal()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	err = ioutil.WriteFile(file, []byte(remoteParsedFileOutput), 0644)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}

func contextShowCommand(c *cli.Context) error {
	fmt.Printf("The current context is: %s\n", kongApi.Config.GetString("current_context"))
	return nil
}

func contextSwitchCommand(c *cli.Context) error {
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

func contextListCommand(c *cli.Context) error {
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
