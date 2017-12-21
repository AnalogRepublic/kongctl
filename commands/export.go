package commands

import (
	"io/ioutil"

	"github.com/analogrepublic/kongctl/data"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var Export = cli.Command{
	Name:    "export",
	Aliases: []string{"download", "fetch", "grab"},
	Usage:   "Export resources to a yaml file",
	Action:  exportCommand,
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

	definition := &data.ServiceDefinition{
		Apis:    remoteApis.Data,
		Plugins: remotePlugins.Data,
	}

	definitionOutput, err := definition.Marshal()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	err = ioutil.WriteFile(file, []byte(definitionOutput), 0644)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}
