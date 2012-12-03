package gosunlight

import (
	"fmt"
	"testing"
)

func TestAPI(t *testing.T) {
	nancy := Legislator{FirstName: "Nancy", LastName: "Pelosi"}
	err := nancy.Get()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(nancy)
}
