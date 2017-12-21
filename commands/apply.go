package commands

import (
	"fmt"

	"github.com/analogrepublic/kongctl/data"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var Apply = cli.Command{
	Name:    "apply",
	Aliases: []string{"upload", "set", "sync"},
	Usage:   "Apply a yaml file to a Kong service",
	Action:  applyCommand,
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
	serviceFile := &data.ServiceDefinitionFile{Name: file}

	// Convert it to a struct
	localServiceDefinition, err := serviceFile.Unmarshal()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	// Grab the remote data & build a service definition from that
	remoteServiceDefinition, err := getRemoteServiceDefinition()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	// Get a diff between what we have locally and what we have remotely
	diff, err := localServiceDefinition.Diff(remoteServiceDefinition)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	fmt.Println(diff)

	// Apply the diff
	applyAdditions(diff)
	applyUpdates(diff)
	applyDeletions(diff)

	return nil
}

func applyUpdates(diff data.ServiceDefinitionDiff) error {
	// Handle Api updates
	for _, updatePair := range diff.Apis.Updates {
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

	return nil
}

func applyAdditions(diff data.ServiceDefinitionDiff) error {
	// Handle Api additions
	for _, api := range diff.Apis.Additions {
		apiAddition := api.(data.Api)

		fmt.Printf("Adding api '%s'\n", apiAddition.Name)

		result, err := kongApi.Apis().Add(&apiAddition)

		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if result.ID == "" {
			return cli.NewExitError(errors.New(fmt.Sprintf("Failed creating api endpoint %s", apiAddition.Name)), 1)
		}
	}

	// Handle Plugin additions
	for _, plugin := range diff.Plugins.Additions {
		pluginAddition := plugin.(data.Plugin)

		fmt.Printf("Adding plugin '%s'\n", pluginAddition.Name)

		result, err := kongApi.Plugins().Add(&pluginAddition)

		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if result.ID == "" {
			return cli.NewExitError(errors.New(fmt.Sprintf("Failed creating plugin endpoint %s", pluginAddition.Name)), 1)
		}
	}

	return nil
}

func applyDeletions(diff data.ServiceDefinitionDiff) error {
	// Handle Api deletions
	for _, deletion := range diff.Apis.Deletions {
		apiDeletion := deletion.(data.Api)
		fmt.Printf("Removing api '%s'\n", apiDeletion.Name)
		err := kongApi.Apis().Delete(&data.ApiRequestParams{Name: apiDeletion.Name})

		if err != nil {
			return cli.NewExitError(errors.New(fmt.Sprintf("Failed removing api endpoint %s", apiDeletion.Name)), 1)
		}
	}

	return nil
}

func getRemoteServiceDefinition() (data.ServiceDefinition, error) {
	var serviceDefinition data.ServiceDefinition

	apis, err := kongApi.Apis().List(&data.ApiRequestParams{})

	if err != nil {
		return serviceDefinition, cli.NewExitError(err, 1)
	}

	plugins, err := kongApi.Plugins().List(&data.PluginRequestParams{})

	if err != nil {
		return serviceDefinition, cli.NewExitError(err, 1)
	}

	// Create a parsed file with the apis we've just got back from
	// the remote kong service
	serviceDefinition = data.ServiceDefinition{
		Apis:    apis.Data,
		Plugins: plugins.Data,
	}

	return serviceDefinition, nil
}
