package kong

import (
	"github.com/analogrepublic/kongctl/data"
	"github.com/dghubble/sling"
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
