package kong

type Plugin struct {
	ID        string
	ApiID     string
	Name      string
	Config    map[string]string
	Enabled   bool
	CreatedAt int
}

type PluginRequestParams struct {
	ID         string
	Name       string
	ApiID      string
	ConsumerID string
	Size       int
	Offset     int
}
