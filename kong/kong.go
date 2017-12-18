package kong

import (
	"net/http"
	"time"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

const (
	// httpTimeout is where we're defining the timeout
	// for the tools HTTP requests.
	httpTimeout = 10 * time.Second
)

// Kong is a type which we'll use
// to interact with the Kong service &
// store any state information about our
// client.
type Kong struct {
	Client *sling.Sling
	Host   string
}

// NewKong should return a new instance of Kong which we
// can use to interact with the API of the service.
func NewKong(host string) (*Kong, error) {
	kong := &Kong{
		Client: newHttpClient(host),
		Host:   host,
	}

	if err := kong.Ping(); err != nil {
		return &Kong{}, errors.Wrap(err, "Test Ping to Kong API failed, unable to reach host.")
	}

	return kong, nil
}

// Plugins returns the plugin handler for the Kong client.
// This should be a handler which can interact with a running
// Kong API.
func (k *Kong) Plugins() *PluginHandler {
	return &PluginHandler{
		Client: k.Client,
		Kong:   k,
	}
}

// Ping makes a single GET request to the base host of
// our Kong service to ensure that the host is reachable.
func (k *Kong) Ping() error {
	_, err := k.Client.Get("/").Request()

	return err
}

// newHttpClient returns a Sling instance which we're going
// to use to wrap the http client logic, this way we only have
// to set the base API url in one place.
func newHttpClient(host string) *sling.Sling {
	client := &http.Client{Timeout: httpTimeout}
	return sling.New().Client(client).Base(host).Set("User-Agent", "Kongctl")
}
