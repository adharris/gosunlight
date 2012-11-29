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

type District struct {
	State  string `json:"state"`
	Number string `json:"number"`
}

func (d District) String() string {
	return fmt.Sprintf("%s-%s", d.State, d.Number)
}

func DistrictsFromZip(zip string) ([]*District, error) {
	var response districtResponse
	p := params{"zip": zip}
	err := districtAPIS.zip.get(&response, p)
	if err != nil {
		return nil, err
	}
	return response.districtSlice(), nil
}

func DistrictsFromLatLong(latitude, longitude float64) ([]*District, error) {
	var response districtResponse
	p := params{
		"latitude":  fmt.Sprintf("%v", latitude),
		"longitude": fmt.Sprintf("%v", longitude),
	}
	err := districtAPIS.latlong.get(&response, p)
	if err != nil {
		return nil, err
	}
	return response.districtSlice(), nil
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
