package data

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type UnparsedYamlFile struct {
	Name     string
	Contents map[interface{}]interface{}
}

// Unmarshal reads a yaml file & unmarshals it into
// a map so we can easily interact with it
func (u *UnparsedYamlFile) Unmarshal() (*ParsedYamlFile, error) {
	parsed := &ParsedYamlFile{}
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

type ParsedYamlFile struct {
	Apis    []*Api
	Plugins []*Plugin
}

// Marshal converts this object into a yaml string, this way we can
// create these "parsed" yaml file objects programatically and then
// output them as a yaml string.
func (p *ParsedYamlFile) Marshal() (string, error) {
	output, err := yaml.Marshal(p)

	if err != nil {
		return "", err
	}

	return string(output), nil
}
