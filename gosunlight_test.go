package gosunlight

import (
	"fmt"
	"testing"
)

func TestAPI(t *testing.T) {
	// match := Legislator{FirstName: "John"}
	reps, err := LegislatorsForZip("02144")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(reps)
	}
}
