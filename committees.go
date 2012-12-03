package gosunlight

import (
	"errors"
	"fmt"
)

var committeeAPIS struct {
	getList       sunlightAPI
	get           sunlightAPI
	forLegislator sunlightAPI
}

func init() {
	committeeAPIS.getList = newSunlightAPI("committees", "getList")
	committeeAPIS.get = newSunlightAPI("committees", "get")
	committeeAPIS.forLegislator = newSunlightAPI("committees", "allForLegislator")
}

// Committee represents a legislative committee from the sunlight api
type Committee struct {
	Chamber       string `json:"chamber"`
	Id            string `json:"id"`
	Name          string `json:"name"`
	Members       []*Legislator
	Subcommittees []*Committee
}

// CommitteeGetList returns a list of committees and their subcommittees
// for a given chamber
//
// See: http://services.sunlightlabs.com/docs/congressapi/committees.getList/
func CommitteeGetList(chamber string) ([]*Committee, error) {
	var response committeesResponse
	p := params{"chamber": chamber}
	err := committeeAPIS.getList.get(&response, p)
	if err != nil {
		return nil, err
	}
	return response.committees(), nil
}

// CommitteeGet returns a committee, its subcommittees, and its members
// given a committee Id
//
// See: http://services.sunlightlabs.com/docs/congressapi/committees.get/
func CommitteeGet(id string) (*Committee, error) {
	var response committeeResponse
	p := params{"id": id}
	err := committeeAPIS.get.get(&response, p)
	if err != nil {
		return nil, err
	}
	return response.committee(), nil
}

// CommitteesForLegislator returns all committees and subcommittees that
// the legislator is a part of.
//
// See: http://services.sunlightlabs.com/docs/congressapi/committees.allForLegislator/
func CommitteesForLegislator(bioguideID string) ([]*Committee, error) {
	var response committeesResponse
	p := params{"bioguide_id": bioguideID}
	err := committeeAPIS.forLegislator.get(&response, p)
	if err != nil {
		return nil, err
	}
	return response.committees(), nil
}

// GetMembers is a convenience wrapper for CommitteeGet which populates
// the Members field of a Committee.
func (c *Committee) GetMembers() error {
	if c.Id == "" {
		return errors.New("Cannot get members of committee: missing id")
	}
	committee, err := CommitteeGet(c.Id)
	if err != nil {
		return err
	}
	c.Members = committee.Members
	return nil
}

type committeeResponse struct {
	Response struct {
		Committee struct {
			Id            string
			Name          string
			Chamber       string
			Subcommittees []struct {
				Committee *Committee
			}
			Members []struct {
				Legislator *Legislator
			}
		}
	}
}

func (cr committeeResponse) committee() *Committee {
	c := Committee{
		Name:          cr.Response.Committee.Name,
		Id:            cr.Response.Committee.Id,
		Chamber:       cr.Response.Committee.Chamber,
		Subcommittees: make([]*Committee, 0, len(cr.Response.Committee.Subcommittees)),
		Members:       make([]*Legislator, 0, len(cr.Response.Committee.Members)),
	}
	for _, m := range cr.Response.Committee.Members {
		c.Members = append(c.Members, m.Legislator)
	}
	for _, sc := range cr.Response.Committee.Subcommittees {
		c.Subcommittees = append(c.Subcommittees, sc.Committee)
	}
	return &c
}

type committeesResponse struct {
	Response struct {
		Committees []struct {
			Committee struct {
				Id            string
				Name          string
				Chamber       string
				Subcommittees []struct {
					Committee *Committee
				}
			}
		}
	}
}

func (lc committeesResponse) committees() []*Committee {
	committees := make([]*Committee, 0, len(lc.Response.Committees))
	for _, c := range lc.Response.Committees {
		committee := &Committee{
			Name:          c.Committee.Name,
			Id:            c.Committee.Id,
			Chamber:       c.Committee.Chamber,
			Subcommittees: make([]*Committee, 0, len(c.Committee.Subcommittees)),
		}
		for _, sc := range c.Committee.Subcommittees {
			committee.Subcommittees = append(committee.Subcommittees, sc.Committee)
		}
		committees = append(committees, committee)
	}
	return committees
}
