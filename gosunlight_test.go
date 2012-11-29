package gosunlight

import (
	"fmt"
	"testing"
)

func TestAPI(t *testing.T) {
	leg := LegislatorSearch("Buch")
	for _, l := range leg {
		fmt.Println(l)
	}

}
