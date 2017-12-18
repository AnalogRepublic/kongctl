package kong

import (
	"github.com/analogrepublic/kongctl/data"
	"github.com/dghubble/sling"
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
func (ph *PluginHandler) List() (*data.PluginList, error) {
	pluginList := &data.PluginList{}

	_, err := ph.Kong.Client.Get(pluginsRootPath).ReceiveSuccess(pluginList)

	if err != nil {
		return pluginList, err
	}

	return pluginList, nil
}
