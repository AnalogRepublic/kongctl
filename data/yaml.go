package data

import (
	"io/ioutil"

	"github.com/analogrepublic/kongctl/util"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// The UnparsedYamlFile struct will represent
// a new file that we need to Parse
type UnparsedYamlFile struct {
	Name string
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

// The ParsedYamlFile represents the struct we'll get
// once we've unmarshaled a yaml file, we can also create
// these when we want to compare what we get from Kong with
// what we've got in files.
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

// ApisIndexed will returbn a list of the Apis indexed by the
// name of the API.
func (p *ParsedYamlFile) ApisIndexed() (map[interface{}]interface{}, error) {
	indexed := map[interface{}]interface{}{}

	for _, api := range p.Apis {
		if api.Name == "" {
			return indexed, errors.New("Missing name for api, unable to index.")
		}

		indexed[api.Name] = *api
	}

	return indexed, nil
}

// Diff should return the differences between the two parsed Yaml
// files. We should get a list of Additions, Deletions & Updates.
func (p *ParsedYamlFile) Diff(other *ParsedYamlFile) (*util.MapDiff, error) {
	// var err error

	// TODO:
	// - Store a map of Plugins by api name or if they're global, in another map (allow for consumers too)

	aApis, err := p.ApisIndexed()

	if err != nil {
		return &util.MapDiff{}, errors.Wrap(err, "Unable to index list of apis in yaml")
	}

	bApis, err := other.ApisIndexed()

	if err != nil {
		return &util.MapDiff{}, errors.Wrap(err, "Unable to index list of apis in Kong")
	}

	// Grab a raw diff of the map, here we're only going
	// to diff the keys, which should be the names of the Apis
	return util.DiffMapsKeys(aApis, bApis), nil
}
