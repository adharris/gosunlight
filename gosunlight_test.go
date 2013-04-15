package gosunlight

import (
	"fmt"
	"testing"
)

func TestAPI(t *testing.T) {
	nancy := Legislator{BioguideID: "B001268"}
	nancy2 := Legislator{BioguideID: "Q000024"}
	l, err := LegislatorGetListAll(&nancy, &nancy2)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(nancy, l)
}
