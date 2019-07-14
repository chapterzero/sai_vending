package machine

import (
	"testing"
)

func TestDisplay(t *testing.T) {
	_ = New(map[Currency]int{C10: 9}, []Inventory{
		Inventory{
			Item{
				Name:  "Canned coffee",
				Price: 120,
			},
			1,
		},
		Inventory{
			Item{
				Name:  "Water PET bottle",
				Price: 100,
			},
			0,
		},
		Inventory{
			Item{
				Name:  "Sport drinks",
				Price: 150,
			},
			1,
		},
	})
}
