package gosunlight

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
)

var legislatorApis struct {
	get     sunlightAPI
	getList sunlightAPI
	search  sunlightAPI
	zip     sunlightAPI
	latlon  sunlightAPI
}

func init() {
	legislatorApis.get = newSunlightAPI("legislators", "get")
	legislatorApis.getList = newSunlightAPI("legislators", "getList")
	legislatorApis.search = newSunlightAPI("legislators", "search")
	legislatorApis.zip = newSunlightAPI("legislators", "allForZip")
	legislatorApis.latlon = newSunlightAPI("legislators", "allForLatLong")
}

// LegislatorSearchThreshold is the threshold used when searching for
// legislators by name using gosunlight.LegislatorSearch().  It should
// be a value from 0 to 1, with 1 being a "perfect match".  Default value
// is .8, values less than .8 are not recommended.
var LegislatorSearchTheshold float64 = .8

// Legislator represents a single legislator from the Sunlight database.
// Available fields are explained here:
//   http://services.sunlightlabs.com/docs/congressapi/legislators.get(List)/
//
// This type is also used for limiting the get(List) functions.
// Create an new instance with just the values to match, and the instance
// as a parameter to the get(List) methods.
type Legislator struct {
	Title            string `json:"title"`
	FirstName        string `json:"firstname"`
	LastName         string `json:"lastname"`
	NameSuffix       string `json:"name_suffix"`
	NickName         string `json:"nickname"`
	Party            string `json:"party"`
	State            string `json:"state"`
	District         string `json:"district"`
	InOffice         bool   `json:"in_office"`
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

	committees []*Committee
}

// String implements fmt.Stringer for Legislators. The Legislator is formatted
// with title, first name, last name, party, and state
//
//   Example: "Sen John Kerry (D MA)"
func (l Legislator) String() string {
	return fmt.Sprintf("%s %s %s (%s %s)", l.Title, l.FirstName, l.LastName, l.Party, l.State)
}

// LegislatorGet gets a single legislator from Sunlight which matches fields
// set in the legislator parameter.
//
// This function only matches legislators who are currently in office.
//
// See: http://services.sunlightlabs.com/docs/congressapi/legislators.get(List)/
func LegislatorGet(legislator Legislator) (*Legislator, error) {
	return getLegislator(false, legislator)
}

// LegislatorGetAll gets a single legislator from Sunlight which matches fields
// set in the legislator parameter.
//
// This function will match against current and past legislators.
//
// See: http://services.sunlightlabs.com/docs/congressapi/legislators.get(List)/
func LegislatorGetAll(legislator Legislator) (*Legislator, error) {
	return getLegislator(true, legislator)
}

func getLegislator(allLegislators bool, legislator Legislator) (*Legislator, error) {
	var l legislatorResponse
	err := legislatorApis.get.get(&l, legislator)
	return l.Response.Legislator, err
}

// LegislatorGetList returns all legislators which match the fields that are
// set in legislators. If multiple legislators are given, they will be
// combined.  If the same field is set on multiple legislator parameters, that
// field will be treated as an OR.
//
// This function only matches legislators who are in office currently.
//
// See http://services.sunlightlabs.com/docs/congressapi/legislators.get(List)/
func LegislatorGetList(legislators ...Legislator) ([]*Legislator, error) {
	return getLegislators(false, legislators...)
}

// LegislatorGetListAll all legislators which match the fields that are set in
// legislators. If multiple legislators are given, they will be combined.
// If the same field is set on multiple legislator parameters, that field will
// be treated as an OR.
//
// This function matches both current and past legislators
//
// See http://services.sunlightlabs.com/docs/congressapi/legislators.get(List)/
func LegislatorGetListAll(legislators ...Legislator) ([]*Legislator, error) {
	return getLegislators(true, legislators...)
}

func getLegislators(allLegislators bool, legislators ...Legislator) ([]*Legislator, error) {
	var r legislatorsResponse
	p := params{}
	if allLegislators {
		p["all_legislators"] = 1
	}
	err := legislatorApis.getList.get(&r, (legislatorSlice)(legislators), p)
	if err != nil {
		return nil, err
	}
	return r.slice(), nil
}

// LegislatorSearch performs a fuzzy search on legislator name.  Each
// legislator is given a score from 0-1, and any legislator above a
// threshold will be returned.  This threshold can be set using the
// package variable LegislatorSearchTheshold.  It defaults to .8, and
// lower values are not recommended.
//
// This function return only legislators who are currently in office.
//
// See: http://services.sunlightlabs.com/docs/congressapi/legislators.search/
func LegislatorSearch(name string) ([]*Legislator, error) {
	return legislatorSearch(name, false)
}

// LegislatorSearchAll performs a fuzzy search on legislator name.  Each
// legislator is given a score from 0-1, and any legislator above a
// threshold will be returned.  This threshold can be set using the
// package variable LegislatorSearchTheshold.  It defaults to .8, and
// lower values are not recommended.
//
// This function matches both current and past legislators.
//
// See: http://services.sunlightlabs.com/docs/congressapi/legislators.search/
func LegislatorSearchAll(name string) ([]*Legislator, error) {
	return legislatorSearch(name, true)
}

func legislatorSearch(name string, allLegislators bool) ([]*Legislator, error) {
	var r legislatorSearchResponse
	p := params{
		"name":            name,
		"threshold":       fmt.Sprintf("%v", LegislatorSearchTheshold),
		"all_legislators": fmt.Sprintf("%v", allLegislators),
	}
	err := legislatorApis.search.get(&r, p)
	if err != nil {
		return nil, err
	}
	return r.slice(), nil
}

// LegislatorsForZip returns all legislators for a 5 digit zip code.
// This function will return 2 senators and at least one house
// representative.  Because zip codes may be in more than one congressional
// district, more than one congressperson may be returned.
//
// This function only returns legislators currently in office.
//
// See: http://services.sunlightlabs.com/docs/congressapi/legislators.allForZip/
func LegislatorsForZip(zip string) ([]*Legislator, error) {
	var r legislatorsResponse
	p := params{"zip": zip}
	err := legislatorApis.zip.get(&r, p)
	if err != nil {
		return nil, err
	}
	return r.slice(), nil
}

// LegislatorsForLatLong returns all legislators for specific latitude and
// longitude.  This usually means one Representative and two Senators.
//
// This function only returns legislators who are currently in office.
//
// See: http://services.sunlightlabs.com/docs/congressapi/legislators.allForLatLong/
func LegislatorsForLatLong(latitude, longitude float64) ([]*Legislator, error) {
	var r legislatorsResponse
	p := params{
		"latitude":  fmt.Sprintf("%v", latitude),
		"longitude": fmt.Sprintf("%v", longitude),
	}
	err := legislatorApis.latlon.get(&r, p)
	if err != nil {
		return nil, err
	}
	return r.slice(), nil
}

// Committees gets a list of the committees and subcommittees that this
// legislator is a part of.  This is a wrapper around CommitteesForLegislator.
// The first call to Committees will block while the committees are fetched
// from sunlight.  Subsequent calls return a cached list.
func (l Legislator) Committees() ([]*Committee, error) {
	if l.committees == nil {
		if l.BioguideID == "" {
			return nil, errors.New("BioguideId missing for legislator.")
		}
		committees, err := CommitteesForLegislator(l.BioguideID)
		if err != nil {
			return nil, err
		}
		l.committees = committees
	}
	return l.committees, nil
}

// Various types used to unmarshal JSON from sunlight.  These
// Types are only used for unmarshaling, individual legislator(s)
// are extracted before returning.
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

//Implementing paramable for a slice of legislators
type legislatorSlice []Legislator

func (ls legislatorSlice) addTo(query *url.Values) {
	for _, l := range ls {
		l.addTo(query)
	}
}

// Implementation of paramable for legislators
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
