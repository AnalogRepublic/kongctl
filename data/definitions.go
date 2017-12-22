package data

import (
	"io/ioutil"

	"github.com/analogrepublic/kongctl/util"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// The ServiceDefinitionFile struct will represent
// a new file that we need to Parse
type ServiceDefinitionFile struct {
	Name string
}

// Unmarshal reads a yaml file & unmarshals it into
// a map so we can easily interact with it
func (u *ServiceDefinitionFile) Unmarshal() (*ServiceDefinition, error) {
	parsed := &ServiceDefinition{}
	contents, err := ioutil.ReadFile(u.Name)

	if err != nil {
		return parsed, err
	}

	err = yaml.Unmarshal(contents, parsed)

	if err != nil {
		return parsed, err
	}

	return parsed, nil
}

// The ServiceDefinition represents the struct we'll get
// once we've unmarshaled a yaml file, we can also create
// these when we want to compare what we get from Kong with
// what we've got in files.
type ServiceDefinition struct {
	Apis    []*Api
	Plugins []*Plugin
}

// The ServiceDefinitionDiff stores a MapDiff for each resource type,
// this should give us an indication of the differences between
// our local config & the Kong service config
type ServiceDefinitionDiff struct {
	Apis    util.MapDiff
	Plugins util.MapDiff
}

// Marshal converts this object into a yaml string, this way we can
// create these "parsed" yaml file objects programatically and then
// output them as a yaml string.
func (p *ServiceDefinition) Marshal() (string, error) {
	output, err := yaml.Marshal(p)

	if err != nil {
		return "", err
	}

	return string(output), nil
}

// ApisIndexed will returbn a list of the Apis indexed by the
// name of the API.
func (p *ServiceDefinition) ApisIndexed() (map[interface{}]interface{}, error) {
	indexed := map[interface{}]interface{}{}

	for _, api := range p.Apis {
		if api.Name == "" {
			return indexed, errors.New("Missing name for api, unable to index.")
		}

		indexed[api.Name] = *api
	}

	return indexed, nil
}

// PluginsIndexed will returbn a list of the Plugins indexed by the
// name of the Plugin.
func (p *ServiceDefinition) PluginsIndexed() (map[interface{}]interface{}, error) {
	indexed := map[interface{}]interface{}{}

	for _, plugin := range p.Plugins {
		if plugin.Name == "" {
			return indexed, errors.New("Missing name for plugin, unable to index.")
		}

		indexed[plugin.Name] = *plugin
	}

	return indexed, nil
}

// Diff should return the differences between the two parsed Yaml
// files. We should get a list of Additions, Deletions & Updates
// for each of the resources.
func (p *ServiceDefinition) Diff(other ServiceDefinition) (ServiceDefinitionDiff, error) {
	var err error
	var apiDiff util.MapDiff
	var pluginDiff util.MapDiff

	if apiDiff, err = p.diffApis(other); err != nil {
		return ServiceDefinitionDiff{}, err
	}

	if pluginDiff, err = p.diffPlugins(other); err != nil {
		return ServiceDefinitionDiff{}, err
	}

	return ServiceDefinitionDiff{
		Apis:    apiDiff,
		Plugins: pluginDiff,
	}, nil
}

func (p *ServiceDefinition) diffApis(other ServiceDefinition) (util.MapDiff, error) {
	var err error
	var diff util.MapDiff
	var a map[interface{}]interface{}
	var b map[interface{}]interface{}

	if a, err = p.ApisIndexed(); err != nil {
		return diff, errors.Wrap(err, "Unable to index list of apis in yaml")
	}

	if b, err = other.ApisIndexed(); err != nil {
		return diff, errors.Wrap(err, "Unable to index list of apis in Kong")
	}

	return *util.DiffMapsKeys(a, b), nil
}

func (p *ServiceDefinition) diffPlugins(other ServiceDefinition) (util.MapDiff, error) {
	var err error
	var diff util.MapDiff
	var a map[interface{}]interface{}
	var b map[interface{}]interface{}

	if a, err = p.PluginsIndexed(); err != nil {
		return diff, errors.Wrap(err, "Unable to index list of plugins in yaml")
	}

	if b, err = other.PluginsIndexed(); err != nil {
		return diff, errors.Wrap(err, "Unable to index list of plugins in Kong")
	}

	return *util.DiffMapsKeys(a, b), nil
}
