package gosunlight

import (
	"fmt"
	"net/url"
	"reflect"
)

var LegislatorSearchTheshold float64 = .8

type Legislator struct {
	Title            string `json:"title"`
	FirstName        string `json:"firstname"`
	LastName         string `json:"lastname"`
	NameSuffix       string `json:"name_suffix"`
	NickName         string `json:"nickname"`
	Party            string `json:"party"`
	State            string `json:"state"`
	District         string `json:"district"`
	InOffice         string `json:"in_office"`
	Gender           string `json:"gender"`
	Phone            string `json:"phone"`
	Fax              string `json:"fax"`
	Website          string `json:"website"`
	WebForm          string `json:"webform"`
	Email            string `json:"email"`
	CongressOffice   string `json:"congress_office"`
	BioguideID       string `json:"bioguide_id""`
	VoteSmartId      string `json:"votesmart_id"`
	FECId            string `json:"fec_id"`
	GovTrackId       string `json:"govtrack_id"`
	CRPID            string `json:"crp_id"`
	CongresspediaURL string `json:"congresspedia_url"`
	TwitterID        string `json:"twitter_id"`
	YouTubeURL       string `json:"youtube_url"`
	FaceBookID       string `json:"facebook_id"`
	SenateClass      string `json:"senate_class"`
	BirthDate        string `json:"birthdate"`
}

func (l Legislator) addTo(query *url.Values) {
	typeOfLegislator := reflect.TypeOf(l)
	valueOfLegislator := reflect.ValueOf(l)

	for i := 0; i < typeOfLegislator.NumField(); i++ {
		jsonKey := typeOfLegislator.Field(i).Tag.Get("json")
		field := valueOfLegislator.Field(i)
		if field.Kind() == reflect.String {
			jsonValue, _ := field.Interface().(string)
			if jsonValue != "" {
				query.Add(jsonKey, jsonValue)
			}
		}
	}
}

type legislatorSlice []Legislator

func (ls legislatorSlice) addTo(query *url.Values) {
	for _, l := range ls {
		l.addTo(query)
	}
}

func (l Legislator) String() string {
	return fmt.Sprintf("%s %s %s (%s %s)", l.Title, l.FirstName, l.LastName, l.Party, l.State)
}

type legislatorResponse struct {
	Response struct {
		Legislator *Legislator
	}
}

type legislatorsResponse struct {
	Response struct {
		Legislators []struct {
			Legislator *Legislator
		}
	}
}

func (lr legislatorsResponse) slice() []*Legislator {
	results := make([]*Legislator, 0, len(lr.Response.Legislators))
	for _, l := range lr.Response.Legislators {
		results = append(results, l.Legislator)
	}
	return results
}

type legislatorSearchResponse struct {
	Response struct {
		Results []struct {
			Result struct {
				Score      float64
				Legislator *Legislator
			}
		}
	}
}

func (lsr legislatorSearchResponse) slice() []*Legislator {
	results := make([]*Legislator, 0, len(lsr.Response.Results))
	for _, l := range lsr.Response.Results {
		results = append(results, l.Result.Legislator)
	}
	return results
}

func LegislatorGet(legislators ...Legislator) *Legislator {
	return getLegislator(false, legislators...)
}

func LegislatorGetAll(legislators ...Legislator) *Legislator {
	return getLegislator(true, legislators...)
}

func getLegislator(allLegislators bool, legislators ...Legislator) *Legislator {
	var l legislatorResponse
	sunlightAPI(&l, "legislators", "get", (legislatorSlice)(legislators))
	return l.Response.Legislator
}

func LegislatorGetList(legislators ...Legislator) []*Legislator {
	return getLegislators(false, legislators...)
}

func LegislatorGetListAll(legislators ...Legislator) []*Legislator {
	return getLegislators(true, legislators...)
}

func getLegislators(allLegislators bool, legislators ...Legislator) []*Legislator {
	var r legislatorsResponse
	p := params{}
	if allLegislators {
		p["all_legislators"] = 1
	}
	sunlightAPI(&r, "legislators", "getList", (legislatorSlice)(legislators), p)
	return r.slice()
}

func LegislatorSearch(name string) []*Legislator {
	return legislatorSearch(name, false)
}

func LegislatorSearchAll(name string) []*Legislator {
	return legislatorSearch(name, true)
}

func legislatorSearch(name string, allLegislators bool) []*Legislator {
	var r legislatorSearchResponse
	p := params{
		"name":            name,
		"threshold":       fmt.Sprintf("%v", LegislatorSearchTheshold),
		"all_legislators": fmt.Sprintf("%v", allLegislators),
	}
	sunlightAPI(&r, "legislators", "search", p)
	return r.slice()
}

func LegislatorsForZip(zip string) []*Legislator {
	var r legislatorsResponse
	p := params{"zip": zip}
	sunlightAPI(&r, "legislators", "allForZip", p)
	return r.slice()
}

func LegislatorsForLatLong(latitude, longitude float64) []*Legislator {
	var r legislatorsResponse
	p := params{
		"latitude":  fmt.Sprintf("%v", latitude),
		"longitude": fmt.Sprintf("%v", longitude),
	}
	sunlightAPI(&r, "legislators", "allForLatLong", p)
	return r.slice()
}
