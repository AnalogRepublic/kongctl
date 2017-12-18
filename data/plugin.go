package data

// Plugin represent the Plugin object we'll
// get back from the Kong API whenever we make a request.
type Plugin struct {
	ID        string
	ApiID     string
	Name      string
	Config    map[string]string
	Enabled   bool
	CreatedAt int
}

// PluginListResponse is an object which represents the
// response we'll get when we're fetching a list of Plugins.
type PluginList struct {
	Total int       `json:"total"`
	Data  []*Plugin `json:"data"`
	Next  string    `json:"next"`
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
