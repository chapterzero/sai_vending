package machine

import (
	"strings"
	"testing"
)

func createTestDisplayMachine() *Machine {
	return New(map[Currency]int{C10: 9}, []Inventory{
		Inventory{
			Item{
				Name:  "Canned coffee",
				Price: 120,
			},
			99,
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
				Name:  "Sport drinks XT",
				Price: 150,
			},
			2,
		},
	})
}

func TestInitialDisplay(t *testing.T) {
	m := createTestDisplayMachine()

	// initial
	expected1 := `
[Input amount]			0 JPY
[Change]				100 JPY			No Change
						10 JPY			Change
[Return gate]			Empty
[Items for sale]
1. Canned coffee		120 JPY
2. Water PET bottle		100 JPY			Sold out
3. Sport drinks XT		150 JPY
[Outlet]				Empty
`

	expected1 = strings.TrimSpace(expected1)

	if expected1 != m.Display() {
		t.Errorf("Expected \n%s\ngot \n%s", expected1, m.Display())
	}

}

func TestDisplayAfterInsert(t *testing.T) {
	m := createTestDisplayMachine()
	m.mainRegister[C100] = 4
	m.Insert(C10)
	m.Insert(C10)
	m.Insert(C100)
	m.Insert(C10)

	expected := `
[Input amount]			130 JPY
[Change]				100 JPY			Change
						10 JPY			Change
[Return gate]			Empty
[Items for sale]
1. Canned coffee		120 JPY			Available for purchase
2. Water PET bottle		100 JPY			Sold out
3. Sport drinks XT		150 JPY
[Outlet]				Empty
`
	expected = strings.TrimSpace(expected)

	if expected != m.Display() {
		t.Errorf("Expected \n%s\ngot \n%s", expected, m.Display())
	}
}

func TestDisplayAfterBuy(t *testing.T) {
	m := createTestDisplayMachine()
	m.Insert(C50)
	m.Insert(C50)
	m.Insert(C100)
	m.Insert(C100)
	m.Buy(0)
	m.Buy(2)

	expected := `
[Input amount]			30 JPY
[Change]				100 JPY			No Change
						10 JPY			No Change
[Return gate]			Empty
[Items for sale]
1. Canned coffee		120 JPY
2. Water PET bottle		100 JPY			Sold out
3. Sport drinks XT		150 JPY
[Outlet]				Canned coffee, Sport drinks XT
`
	expected = strings.TrimSpace(expected)

	if expected != m.Display() {
		t.Errorf("Expected \n%s\ngot \n%s", expected, m.Display())
	}
}

func TestDisplayAfterReturnAndGetItems(t *testing.T) {
	m := createTestDisplayMachine()
	m.Insert(C50)
	m.Insert(C50)
	m.Insert(C100)
	m.Insert(C100)
	m.Buy(0)
	m.Buy(2)
	m.ReturnInput()
	m.GetItems()

	expected := `
[Input amount]			0 JPY
[Change]				100 JPY			No Change
						10 JPY			No Change
[Return gate]			10 JPY, 10 JPY, 10 JPY
[Items for sale]
1. Canned coffee		120 JPY
2. Water PET bottle		100 JPY			Sold out
3. Sport drinks XT		150 JPY
[Outlet]				Empty
`
	expected = strings.TrimSpace(expected)

	if expected != m.Display() {
		t.Errorf("Expected \n%s\ngot \n%s", expected, m.Display())
	}
}

func TestDisplayAfterGetReturn(t *testing.T) {
	m := createTestDisplayMachine()
	m.Insert(C50)
	m.Insert(C50)
	m.Insert(C100)
	m.Insert(C100)
	m.Buy(0)
	m.Buy(2)
	m.ReturnInput()
	m.GetItems()
	m.GetReturn()

	expected := `
[Input amount]			0 JPY
[Change]				100 JPY			No Change
						10 JPY			No Change
[Return gate]			Empty
[Items for sale]
1. Canned coffee		120 JPY
2. Water PET bottle		100 JPY			Sold out
3. Sport drinks XT		150 JPY
[Outlet]				Empty
`
	expected = strings.TrimSpace(expected)

	if expected != m.Display() {
		t.Errorf("Expected \n%s\ngot \n%s", expected, m.Display())
	}
}
