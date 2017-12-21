package data

import (
	"errors"
)

// Api represent the Api object we'll
// get back from the Kong API whenever we make a request.
type Api struct {
	ID                     string      `json:"id,omitempty" yaml:"-"`
	Name                   string      `json:"name,omitempty" yaml:"name,omitempty"`
	Uris                   interface{} `json:"uris,omitempty" yaml:"uris,omitempty"`
	Hosts                  interface{} `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	Methods                interface{} `json:"methods,omitempty" yaml:"methods,omitempty"`
	PreserveHost           bool        `json:"preserve_host" yaml:"preserve_host,omitempty"`
	StripUri               bool        `json:"strip_uri" yaml:"strip_uri,omitempty"`
	UpstreamUrl            string      `json:"upstream_url,omitempty" yaml:"upstream_url,omitempty"`
	UpstreamConnectTimeout int         `json:"upstream_connect_timeout,omitempty" yaml:"upstream_connect_timeout,omitempty"`
	UpstreamReadTimeout    int         `json:"upstream_read_timeout,omitempty" yaml:"upstream_read_timeout,omitempty"`
	UpstreamSendTimeout    int         `json:"upstream_send_timeout,omitempty" yaml:"upstream_send_timeout,omitempty"`
	HttpIfTerminated       bool        `json:"http_if_terminated" yaml:"http_if_terminated,omitempty"`
	HttpsOnly              bool        `json:"https_only" yaml:"https_only,omitempty"`
	Retries                int         `json:"retries,omitempty" yaml:"retries,omitempty"`
	CreatedAt              int         `json:"created_at,omitempty" yaml:"-"`
}

// ApiList is an object which represents the
// response we'll get when we're fetching a list of Apis.
type ApiList struct {
	Total int    `json:"total"`
	Data  []*Api `json:"data"`
	Next  string `json:"next"`
}

// ApiRequestParams allows us to pass in a query
// string of parameters to some of the Api requests.
type ApiRequestParams struct {
	ID          string `url:"id,omitempty"`
	Name        string `url:"name,omitempty"`
	UpstreamUrl string `url:"upstream_url,omitempty"`
	Retries     int    `url:"retries,omitempty"`
	Size        int    `url:"size_id,omitempty"`
	Offset      int    `url:"offset_id,omitempty"`
}

// Identifier should grab the identifier we've passed into
// our request params, favouring the ID over the name.
func (arp *ApiRequestParams) Identifier() (string, error) {
	if arp.ID != "" {
		return arp.ID, nil
	}

	if arp.Name != "" {
		return arp.Name, nil
	}

	return "", errors.New("You must provide an ID or Name in the ApiRequestParams")
}
