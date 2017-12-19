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

// Add will create a new api resource on Kong & Handle any conflicts.
func (ah *ApiHandler) Add(api *data.Api) (*data.Api, error) {
	respApi := &data.Api{}

	_, err := ah.Kong.Client.Post(apisRootPath).BodyJSON(api).ReceiveSuccess(respApi)

	if err != nil {
		return respApi, err
	}

	return respApi, nil
}

// Update will make a PUT request to update an existing
// api stored in the Kong service
func (ah *ApiHandler) Update(params *data.ApiRequestParams, updatedData *data.Api) (*data.Api, error) {
	respApi := &data.Api{}

	identifier, err := params.Identifier()

	if err != nil {
		return respApi, errors.Wrap(err, "You must provide an ID or Name to update an Api")
	}

	path := fmt.Sprintf("%s/%s", apisRootPath, identifier)

	_, err = ah.Kong.Client.Patch(path).BodyJSON(updatedData).ReceiveSuccess(respApi)

	if err != nil {
		return respApi, err
	}

	return respApi, nil
}

// Delete will make a DELETE request to remove an api from the Kong service
func (ah *ApiHandler) Delete(params *data.ApiRequestParams) error {
	identifier, err := params.Identifier()

	if err != nil {
		return errors.Wrap(err, "You must provide an ID or Name to update an Api")
	}

	path := fmt.Sprintf("%s/%s", apisRootPath, identifier)

	request, err := ah.Kong.Client.Delete(path).Request()

	if err != nil {
		return err
	}

	_, err = ah.Kong.Client.Do(request, nil, nil)

	if err != nil {
		return err
	}

	return nil
}
