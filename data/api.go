package data

// Api represent the Api object we'll
// get back from the Kong API whenever we make a request.
type Api struct {
	ID                     string      `json:"id,omitempty"`
	Hosts                  interface{} `json:"hosts,omitempty"`
	Methods                interface{} `json:"methods,omitempty"`
	Uris                   interface{} `json:"uris,omitempty"`
	Name                   string      `json:"name,omitempty"`
	HttpIfTerminated       bool        `json:"http_if_terminated,omitempty"`
	HttpsOnly              bool        `json:"https_only,omitempty"`
	PreserveHost           bool        `json:"preserve_host,omitempty"`
	StripUri               bool        `json:"strip_uri,omitempty"`
	UpstreamConnectTimeout int         `json:"upstream_connect_timeout,omitempty"`
	UpstreamReadTimeout    int         `json:"upstream_read_timeout,omitempty"`
	UpstreamSendTimeout    int         `json:"upstream_send_timeout,omitempty"`
	UpstreamUrl            string      `json:"upstream_url,omitempty"`
	Retries                int         `json:"retries,omitempty"`
	CreatedAt              int         `json:"created_at,omitempty"`
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
