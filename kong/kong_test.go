package kong

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/analogrepublic/kongctl/config"
	"github.com/analogrepublic/kongctl/data"
	"github.com/stretchr/testify/assert"
)

var conf *config.Config
var seed = strconv.Itoa(int(time.Now().Unix()))

var testApiObject *data.Api

func getConfig() *config.Config {
	if conf == nil {
		err := config.Init()

		if err != nil {
			panic(fmt.Sprintf("Unable to load configuration: %s", err))
		}

		conf = config.GetConfig()
	}

	return conf
}

func getKong() (*Kong, error) {
	c := getConfig()
	context, err := c.GetCurrentContext()

	if err != nil {
		panic("Unable to determine context to use, please check your config file")
	}

	return NewKong(context.Host, c)
}

func TestNewKong(t *testing.T) {
	kongApi, err := getKong()

	assert.Nil(t, err)
	assert.NotNil(t, kongApi)
}

func TestApisAdd(t *testing.T) {
	kongApi, err := getKong()

	assert.Nil(t, err)

	apiItem := &data.Api{
		Name: "Testing" + seed,
		Uris: []string{
			"/testing" + seed,
		},
		StripUri:    true,
		UpstreamUrl: "https://google.com",
	}

	api, err := kongApi.Apis().Add(apiItem)

	assert.Nil(t, err)
	assert.NotEqual(t, api.ID, "")
	assert.Equal(t, api.Name, apiItem.Name)

	testApiObject = api
}

func TestApisList(t *testing.T) {
	kongApi, err := getKong()

	assert.Nil(t, err)

	apis, err := kongApi.Apis().List(&data.ApiRequestParams{})

	assert.Nil(t, err)
	assert.NotEqual(t, len(apis.Data), 0)
}

func TestApisRetrieve(t *testing.T) {
	kongApi, err := getKong()

	assert.Nil(t, err)

	api, err := kongApi.Apis().Retrieve(&data.ApiRequestParams{Name: testApiObject.Name})

	assert.Nil(t, err)
	assert.NotEqual(t, api.ID, "")
	assert.Equal(t, api.Name, testApiObject.Name)
}

func TestApisUpdate(t *testing.T) {
	kongApi, err := getKong()

	assert.Nil(t, err)

	updatedData := &data.Api{
		Uris: []string{
			"/testing-updated",
		},
		UpstreamUrl: "https://google.com/updated",
	}

	api, err := kongApi.Apis().Update(&data.ApiRequestParams{Name: testApiObject.Name}, updatedData)

	assert.Nil(t, err)
	assert.NotEqual(t, api.ID, "")
	assert.Equal(t, api.UpstreamUrl, updatedData.UpstreamUrl)
}

func TestApisDelete(t *testing.T) {
	kongApi, err := getKong()

	assert.Nil(t, err)

	err = kongApi.Apis().Delete(&data.ApiRequestParams{Name: testApiObject.Name})

	assert.Nil(t, err)
}

// func TestPluginsList(t *testing.T) {
// kongApi, err := getKong()

//  assert.Nil(t, err)

//  plugins, err := kongApi.Plugins().List(&data.PluginRequestParams{})

//  assert.Nil(t, err)
//  assert.NotEqual(t, len(plugins.Data), 0)
// }
