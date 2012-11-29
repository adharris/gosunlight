// Package gosunlight is an implementation of the Sunlight labs API
// See http://services.sunlightlabs.com/
//
// To use any of the functions included, a SunlightLabs API key must
// provided.  API keys can be acquired at
// http://services.sunlightlabs.com/accounts/register/
//
// To provide the Sunlight API Key to go sunlight, simply set
// gosunlight.SunlightKey to your API key.  If no key is provided,
// gosunlight will attempt to load the sunlight key from the OS
// environment variable SUNLIGHT_KEY.
package gosunlight

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const (
	sunlightURL = "http://services.sunlightlabs.com/api/"
)

// The API Key for Sunlight Labs.  This can be set manually, or it
// will be pulled from the SUNLIGHT_KEY environment variable.
var SunlightKey string

// Pull the api key from the environment variables if it has not
// been set in code.
func init() {
	osKey := os.Getenv("SUNLIGHT_KEY")
	if osKey != "" && SunlightKey == "" {
		SunlightKey = osKey
	}
}

// An interface for types that can be translated to url parameters
type paramable interface {
	// adds the parameters in this type to a url.Values object
	addTo(query *url.Values)
}

// A simple map that implements the paramable interface
type params map[string]interface{}

// Implements paramable. Adds all values from the map to the query.
func (p params) addTo(query *url.Values) {
	for key := range p {
		query.Add(key, fmt.Sprintf("%v", p[key]))
	}
}

// A type for a specific api call
type sunlightAPI struct {
	api    string
	method string
	rawURL string
}

// Returns a sunlight api handler
func newSunlightAPI(api, method string) sunlightAPI {
	return sunlightAPI{
		api:    api,
		method: method,
		rawURL: fmt.Sprintf("%s%s.%s.json", sunlightURL, api, method),
	}
}

// Runs the api request.  The JSON response is unmarshaled into
// the v parameter
func (api sunlightAPI) get(v interface{}, params ...paramable) error {

	fullURL, _ := url.Parse(api.rawURL)
	query := fullURL.Query()
	query.Add("apikey", SunlightKey)
	for _, p := range params {
		p.addTo(&query)
	}
	fullURL.RawQuery = query.Encode()

	res, err := http.Get(fullURL.String())
	if err != nil {
		return err
	}

	if res.StatusCode == 400 {
		errorMessage, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(errorMessage))
	} else {
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&v)
		if err != nil {
			return err
		}
	}
	return nil

}
