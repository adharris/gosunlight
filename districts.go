package gosunlight

import (
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
}

// String implements fmt.Stringer for Districts
func (d District) String() string {
	return fmt.Sprintf("%s-%s", d.State, d.Number)
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
