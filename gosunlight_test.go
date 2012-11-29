package gosunlight

import (
	"fmt"
	"testing"
)

func TestAPI(t *testing.T) {
	districts, _ := DistrictsFromZip("02144")
	for _, d := range districts {
		fmt.Println(d)
		rep, _ := d.Representative()
		fmt.Printf("  %v\n", rep)
		sens, _ := d.Sentators()
		for _, s := range sens {
			fmt.Printf("  %v\n", s)
		}
	}
}
