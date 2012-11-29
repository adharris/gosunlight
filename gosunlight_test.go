package gosunlight

import (
	"fmt"
	"testing"
)

func TestAPI(t *testing.T) {
	districts, err := DistrictsFromLatLong(35.778788, -78.787805)
	fmt.Println(err)
	for _, d := range districts {
		fmt.Println(d)
	}
}
