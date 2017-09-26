package kong

import (
	"net/http"
	"time"
)

const (
	httpTimeout = 10 * time.Second
)

type Kong struct {
	Client *http.Client
	Host   string
}

func NewKong(host string, errors chan error) *Kong {
	kong := &Kong{
		Client: newHttpClient(),
		Host:   host,
	}

	_, err := kong.Client.Get(kong.Host)

	if err != nil {
		errors <- err
	}

	return kong
}

func newHttpClient() *http.Client {
	return &http.Client{
		Timeout: httpTimeout,
	}
}

func (k *Kong) Plugins() ([]*Plugin, error) {
	resp, err := k.Client.Get(k.Host + "/plugins")
}
