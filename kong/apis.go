package kong

import (
	"fmt"

	"github.com/analogrepublic/kongctl/data"
	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

const (
	// apisRootPath is explicitly set for this package.
	apisRootPath = "/apis"
)

// ApiHandler is our object to interface with the Plugin
// side of the kong API.
type ApiHandler struct {
	Client *sling.Sling
	Kong   *Kong
}

// List will make a GET request with our request params and
// return a ApiList which contains the number of apis,
// a list of the apis fetched and a reference to the next page.
func (ah *ApiHandler) List(params *data.ApiRequestParams) (*data.ApiList, error) {
	apiList := &data.ApiList{}

	_, err := ah.Kong.Client.Get(apisRootPath).QueryStruct(params).ReceiveSuccess(apiList)

	if err != nil {
		return apiList, err
	}

	return apiList, nil
}

// Retrieve will make a GET request to fetch a single Api by Name or
// ID which will be provided in the params
func (ah *ApiHandler) Retrieve(params *data.ApiRequestParams) (*data.Api, error) {
	api := &data.Api{}
	identifier, err := params.Identifier()

	if err != nil {
		return api, errors.Wrap(err, "You must provide an ID or Name to retrieve an Api")
	}

	path := fmt.Sprintf("%s/%s", apisRootPath, identifier)

	_, err = ah.Kong.Client.Get(path).ReceiveSuccess(api)

	if err != nil {
		return api, err
	}

	return api, nil
}
