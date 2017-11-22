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

func NewKong(host string) (*Kong, error) {
	kong := &Kong{
		Client: newHttpClient(),
		Host:   host,
	}

	_, err := kong.Client.Get(kong.Host)

	if err != nil {
		return &Kong{}, err
	}

	return kong, nil
}

func newHttpClient() *http.Client {
	return &http.Client{
		Timeout: httpTimeout,
	}
}

func (k *Kong) Plugins() ([]*Plugin, error) {
	plugins := []*Plugin{}
	resp, err := k.Client.Get(k.Host + "/plugins")

	if err != nil {
		return plugins, err
	}

	return plugins, nil
}
