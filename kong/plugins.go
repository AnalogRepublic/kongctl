package kong

import (
	"fmt"

	"github.com/analogrepublic/kongctl/data"
	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

const (
	// pluginsRootPath is explicitly set for this package.
	pluginsRootPath = "/plugins"
)

// PluginHandler is our object to interface with the Plugin
// side of the kong API.
type PluginHandler struct {
	Client *sling.Sling
	Kong   *Kong
}

// List will make a GET request with our request params and
// return a PluginList which contains the number of plugins,
// a list of the plugins fetched and a reference to the next page.
func (ph *PluginHandler) List(params *data.PluginRequestParams) (*data.PluginList, error) {
	pluginList := &data.PluginList{}

	_, err := ph.Kong.Client.Get(pluginsRootPath).QueryStruct(params).ReceiveSuccess(pluginList)

	if err != nil {
		return pluginList, err
	}

	return pluginList, nil
}

// Retrieve will make a GET request to fetch a single plugin by Name or
// ID which will be provided in the params
func (ph *PluginHandler) Retrieve(params *data.PluginRequestParams) (*data.Plugin, error) {
	plugin := &data.Plugin{}
	identifier, err := params.Identifier()

	if err != nil {
		return plugin, errors.Wrap(err, "You must provide an ID or Name to retrieve a plugin")
	}

	path := fmt.Sprintf("%s/%s", pluginsRootPath, identifier)

	_, err = ph.Kong.Client.Get(path).ReceiveSuccess(plugin)

	if err != nil {
		return plugin, err
	}

	return plugin, nil
}

// Add will create a new plugin resource on Kong & Handle any conflicts.
func (ph *PluginHandler) Add(plugin *data.Plugin) (*data.Plugin, error) {
	respPlugin := &data.Plugin{}

	// By default we'll be creating a global plugin
	path := pluginsRootPath

	// If we're only applying it to an api, use the apis
	// plugin endpoint
	if plugin.ApiID != "" && plugin.ConsumerID == "" {
		path = fmt.Sprintf("%s/%s/%s", apisRootPath, plugin.ConsumerID, path)
	}

	_, err := ph.Kong.Client.Post(path).BodyJSON(plugin).ReceiveSuccess(respPlugin)

	if err != nil {
		return respPlugin, err
	}

	return respPlugin, nil
}

// Update will make a PUT request to update an existing
// plugin stored in the Kong service
func (ph *PluginHandler) Update(params *data.PluginRequestParams, updatedData *data.Plugin) (*data.Plugin, error) {
	respPlugin := &data.Plugin{}

	identifier, err := params.Identifier()

	if err != nil {
		return respPlugin, errors.Wrap(err, "You must provide an ID or Name to update a plugin")
	}

	path := fmt.Sprintf("%s/%s", pluginsRootPath, identifier)

	_, err = ph.Kong.Client.Patch(path).BodyJSON(updatedData).ReceiveSuccess(respPlugin)

	if err != nil {
		return respPlugin, err
	}

	return respPlugin, nil
}

// Delete will make a DELETE request to remove a plugin from the Kong service
func (ph *PluginHandler) Delete(params *data.PluginRequestParams) error {
	identifier, err := params.Identifier()

	if err != nil {
		return errors.Wrap(err, "You must provide an ID or Name to delete a plugin")
	}

	path := fmt.Sprintf("%s/%s", pluginsRootPath, identifier)

	request, err := ph.Kong.Client.Delete(path).Request()

	if err != nil {
		return err
	}

	_, err = ph.Kong.Client.Do(request, nil, nil)

	if err != nil {
		return err
	}

	return nil
}
