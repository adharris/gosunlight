package gosunlight

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

const (
	sunlightURL = "http://services.sunlightlabs.com/api/"
)

var sunlightKey string

func init() {
	sunlightKey = os.Getenv("SUNLIGHT_KEY")
}

type paramable interface {
	addTo(query *url.Values)
}

type params map[string]interface{}

func (p params) addTo(query *url.Values) {
	for key := range p {
		query.Add(key, fmt.Sprintf("%v", p[key]))
	}
}

func sunlightAPI(v interface{}, api, method string, params ...paramable) error {
	baseURL := fmt.Sprintf("%s%s.%s.json?apikey=%s", sunlightURL, api, method, sunlightKey)

	fullURL, _ := url.Parse(baseURL)
	query := fullURL.Query()
	for _, p := range params {
		p.addTo(&query)
	}
	fullURL.RawQuery = query.Encode()

	res, err := http.Get(fullURL.String())
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&v)
	if err != nil {
		return err
	}
	return nil
}
