package data

// Plugin represent the Plugin object we'll
// get back from the Kong API whenever we make a request.
type Plugin struct {
	ID        string      `json:"id,omitempty" yaml:"id,omitempty"`
	ApiID     string      `json:"api_id,omitempty" yaml:"api_id,omitempty"`
	Name      string      `json:"name,omitempty" yaml:"name,omitempty"`
	Config    interface{} `json:"config,omitempty" yaml:"config,omitempty"`
	Enabled   bool        `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	CreatedAt int         `json:"created_at,omitempty" yaml:"created_at,omitempty"`
}

// PluginList is an object which represents the
// response we'll get when we're fetching a list of Plugins.
type PluginList struct {
	Total int       `json:"total"`
	Data  []*Plugin `json:"data"`
	Next  string    `json:"next"`
}

// FilterData allows you to provide a filter callback to fetch
// a new filtered sub-set of the Data
func (pl *PluginList) FilterData(f func(Plugin) bool) []Plugin {
	newList := make([]Plugin, 0)

	for _, plugin := range pl.Data {
		if f(*plugin) {
			newList = append(newList, *plugin)
		}
	}

	return newList
}

// PluginRequestParams allows us to pass in a query
// string of parameters to some of the plugin requests.
type PluginRequestParams struct {
	ID         string `url:"id,omitempty"`
	Name       string `url:"name,omitempty"`
	ApiID      string `url:"api_id,omitempty"`
	ConsumerID string `url:"consumer_id,omitempty"`
	Size       int    `url:"size_id,omitempty"`
	Offset     int    `url:"offset_id,omitempty"`
}
