package kong

import (
	"strconv"
	"testing"
	"time"

	"github.com/analogrepublic/kongctl/data"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var seed = strconv.Itoa(int(time.Now().Unix()))

var testApiObject *data.Api

func config() *viper.Viper {
	v := viper.New()

	v.Set("host", "http://testing.kongctl.io:8001")

	return v
}

func TestNewKong(t *testing.T) {
	c := config()

	kong, err := NewKong(c.GetString("host"), c)

	assert.Nil(t, err)
	assert.NotNil(t, kong)
}

func TestApisAdd(t *testing.T) {
	c := config()

	kong, err := NewKong(c.GetString("host"), c)

	assert.Nil(t, err)

	apiItem := &data.Api{
		Name: "Testing" + seed,
		Uris: []string{
			"/testing" + seed,
		},
		StripUri:    true,
		UpstreamUrl: "https://google.com",
	}

	api, err := kong.Apis().Add(apiItem)

	assert.Nil(t, err)
	assert.NotEqual(t, api.ID, "")
	assert.Equal(t, api.Name, apiItem.Name)

	testApiObject = api
}

func TestApisList(t *testing.T) {
	c := config()

	kong, err := NewKong(c.GetString("host"), c)

	assert.Nil(t, err)

	apis, err := kong.Apis().List(&data.ApiRequestParams{})

	assert.Nil(t, err)
	assert.NotEqual(t, len(apis.Data), 0)
}

func TestApisRetrieve(t *testing.T) {
	c := config()

	kong, err := NewKong(c.GetString("host"), c)

	assert.Nil(t, err)

	api, err := kong.Apis().Retrieve(&data.ApiRequestParams{Name: testApiObject.Name})

	assert.Nil(t, err)
	assert.NotEqual(t, api.ID, "")
	assert.Equal(t, api.Name, testApiObject.Name)
}

func TestApisUpdate(t *testing.T) {
	c := config()

	kong, err := NewKong(c.GetString("host"), c)

	assert.Nil(t, err)

	updatedData := &data.Api{
		Uris: []string{
			"/testing-updated",
		},
		UpstreamUrl: "https://google.com/updated",
	}

	api, err := kong.Apis().Update(&data.ApiRequestParams{Name: testApiObject.Name}, updatedData)

	assert.Nil(t, err)
	assert.NotEqual(t, api.ID, "")
	assert.Equal(t, api.UpstreamUrl, updatedData.UpstreamUrl)
}

func TestApisDelete(t *testing.T) {
	c := config()

	kong, err := NewKong(c.GetString("host"), c)

	assert.Nil(t, err)

	err = kong.Apis().Delete(&data.ApiRequestParams{Name: testApiObject.Name})

	assert.Nil(t, err)
}

// func TestPluginsList(t *testing.T) {
//  c := config()

//  kong, err := NewKong(c.GetString("host"), c)

//  assert.Nil(t, err)

//  plugins, err := kong.Plugins().List(&data.PluginRequestParams{})

//  assert.Nil(t, err)
//  assert.NotEqual(t, len(plugins.Data), 0)
// }
