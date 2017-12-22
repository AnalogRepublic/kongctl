package kong

import (
	"net/http"
	"time"

	"github.com/analogrepublic/kongctl/config"
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
	Config *config.Config
	Client *sling.Sling
	Host   string
}

type ApiError struct {
	Message string `json:"message"`
}

// NewKong should return a new instance of Kong which we
// can use to interact with the API of the service.
func NewKong(host string, c *config.Config) (*Kong, error) {
	kong := &Kong{
		Config: c,
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

// Apis returns the api handler for the Kong client.
// This should be a handler which can interact with a running
// Kong API.
func (k *Kong) Apis() *ApiHandler {
	return &ApiHandler{
		Client: k.Client,
		Kong:   k,
	}
}

// Ping makes a single GET request to the base host of
// our Kong service to ensure that the host is reachable.
func (k *Kong) Ping() error {
	req, err := k.Client.Get("/").Request()

	if err != nil {
		return err
	}

	var success interface{}
	var fail interface{}

	_, err = k.Client.Do(req, success, fail)

	if err != nil {
		return err
	}

	if fail != nil {
		return errors.New("The request did not return with a 2XX response")
	}

	return nil
}

// newHttpClient returns a Sling instance which we're going
// to use to wrap the http client logic, this way we only have
// to set the base API url in one place.
func newHttpClient(host string) *sling.Sling {
	client := &http.Client{Timeout: httpTimeout}
	return sling.New().Client(client).Base(host).Set("User-Agent", "Kongctl")
}
