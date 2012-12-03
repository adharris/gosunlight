package gosunlight

import (
	"errors"
	"fmt"
)

var districtAPIS struct {
	zip     sunlightAPI
	latlong sunlightAPI
}

func init() {
	districtAPIS.zip = newSunlightAPI("districts", "getDistrictsFromZip")
	districtAPIS.latlong = newSunlightAPI("districts", "getDistrictFromLatLong")
}

// District represents a congressional district.
type District struct {
	State  string `json:"state"`
	Number string `json:"number"`

	rep      *Legislator
	senators []*Legislator
}

// DistrictsFromZip returns a list of districts for a given zip code.  Because
// a zip code may be in more than one district, this function may return
// more than once district.
//
// See: http://services.sunlightlabs.com/docs/congressapi/districts.getDistrictsFromZip/
func DistrictsFromZip(zip string) ([]*District, error) {
	var response districtResponse
	p := params{"zip": zip}
	err := districtAPIS.zip.get(&response, p)
	if err != nil {
		return nil, err
	}
	return response.districtSlice(), nil
}

// DistrictsFromLatLong returns a single district for a given latitude and
// longitude.
//
// See: http://services.sunlightlabs.com/docs/congressapi/districts.getDistrictFromLatLong/
func DistrictFromLatLong(latitude, longitude float64) (*District, error) {
	return districtFromLatLong(latitude, longitude, 2010)
}

// DistrictsFromLatLong2012 returns the district for a point based on the
// 2012 redistricting.
//
// See: http://services.sunlightlabs.com/docs/congressapi/districts.getDistrictFromLatLong/
func DistrictFromLatLong2012(latitude, longitude float64) (*District, error) {
	return districtFromLatLong(latitude, longitude, 2012)
}

func districtFromLatLong(latitude, longitude float64, districts int) (*District, error) {
	var response districtResponse
	p := params{
		"latitude":  fmt.Sprintf("%v", latitude),
		"longitude": fmt.Sprintf("%v", longitude),
		"districts": districts,
	}
	err := districtAPIS.latlong.get(&response, p)
	if err != nil {
		return nil, err
	}
	return response.districtSlice()[0], nil
}

// Representative returns the house of representatives member for a given
// district.  This function will block while the data is fetched from
// sunlight.  Subsequent calls return a cached value.
func (d *District) Representative() (*Legislator, error) {
	if d.rep == nil {
		if d.State == "" || d.Number == "" {
			return nil, errors.New("State or number missing from district; cannot get legislators")
		}
		legislator, err := LegislatorGet(Legislator{State: d.State, District: d.Number})
		if err != nil {
			return nil, err
		}
		d.rep = legislator
	}
	return d.rep, nil
}

// Senators return the senators for a given district.  This function will
// block while the data is fetched from sunlight.  Subsequent calls return
// a cached value.
func (d *District) Sentators() ([]*Legislator, error) {
	if d.senators == nil {
		if d.State == "" {
			return nil, errors.New("State missing from district; cannot get senators")
		}
		legislators, err := LegislatorGetList(Legislator{Title: "Sen", State: d.State})
		if err != nil {
			return nil, err
		}
		d.senators = legislators
	}
	return d.senators, nil
}

type districtResponse struct {
	Response struct {
		Districts []struct {
			District *District
		}
	}
}

func (dr districtResponse) districtSlice() []*District {
	response := make([]*District, 0, len(dr.Response.Districts))
	for _, d := range dr.Response.Districts {
		response = append(response, d.District)
	}
	return response
}
