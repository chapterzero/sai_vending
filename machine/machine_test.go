package machine

import (
	"testing"
)

func TestMachineInsert(t *testing.T) {
	testCases := []struct {
		name          string
		m             *Machine
		input         Currency
		expectedError string
	}{
		{
			name:          "Inserting 10 coin on empty machine",
			m:             createEmptyMachine(),
			input:         C10,
			expectedError: "",
		},
		{
			name:          "Inserting 50 coin on empty machine",
			m:             createEmptyMachine(),
			input:         C50,
			expectedError: "Unable to return change",
		},
		{
			name:          "Inserting 50 coin on machine running out of 10 coins",
			m:             New(map[Currency]int{C10: 5}, []Inventory{}),
			input:         C50,
			expectedError: "Unable to return change",
		},
		{
			name:          "Inserting 100 coin on machine running out of 10 coins",
			m:             New(map[Currency]int{C10: 5}, []Inventory{}),
			input:         C100,
			expectedError: "Unable to return change",
		},
		{
			name:          "Successfully Inserting 100 coin on machine",
			m:             New(map[Currency]int{C10: 10}, []Inventory{}),
			input:         C100,
			expectedError: "",
		},
		{
			name:          "Inserting 500 coin on machine running out of 100 coins",
			m:             New(map[Currency]int{C10: 20, C100: 3}, []Inventory{}),
			input:         C500,
			expectedError: "Unable to return change",
		},
		{
			name:          "Successfully Inserting 500 coin on machine",
			m:             New(map[Currency]int{C10: 20, C100: 5}, []Inventory{}),
			input:         C500,
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.m.Insert(tc.input)
			if tc.expectedError == "" {
				if err != nil {
					t.Errorf("Expected no error occured, got %s", err)
				}
			} else {
				if err.Error() != tc.expectedError {
					t.Errorf("Expected error '%s', got '%s'", tc.expectedError, err.Error())
				}
			}
		})
	}
}

func TestMachineBuyShouldReturnErrorOnInvalidIdx(t *testing.T) {
	m := &Machine{
		inputRegister: []Currency{500},
		inventories: []Inventory{
			Inventory{
				Item{
					Name:  "Item 1",
					Price: 100,
				},
				99,
			},
			Inventory{
				Item{
					Name:  "Item 2",
					Price: 110,
				},
				99,
			},
		},
	}

	err := m.Buy(-1)
	if err == nil {
		t.Errorf("Expected error not nil, got nil")
	}

	err = m.Buy(2)
	if err == nil {
		t.Errorf("Expected error not nil, got nil")
	}
}

func TestMachineBuyInsertedMoneyNotEnough(t *testing.T) {
	m := &Machine{
		mainRegister:  map[Currency]int{C10: 9},
		inputRegister: []Currency{C100},
		inventories: []Inventory{
			Inventory{
				Item{
					Name:  "Item 1",
					Price: 100,
				},
				99,
			},
			Inventory{
				Item{
					Name:  "Item 2",
					Price: 110,
				},
				99,
			},
		},
	}

	err := m.Buy(1)
	if err == nil {
		t.Errorf("Expected error not nil, got nil")
	}
}

func TestMachineBuySuccessful(t *testing.T) {
	m := &Machine{
		mainRegister:  map[Currency]int{C10: 9},
		inputRegister: []Currency{C100, C100},
		inventories: []Inventory{
			Inventory{
				Item{
					Name:  "Item 1",
					Price: 100,
				},
				99,
			},
			Inventory{
				Item{
					Name:  "Item 2",
					Price: 110,
				},
				99,
			},
		},
	}

	err := m.Buy(1)
	if err != nil {
		t.Errorf("Expected error nil, got %v", err.Error())
	}

	if len(m.inputRegister) != 9 {
		t.Errorf("Expected input register contain 9 coins, got %d", len(m.inputRegister))
	}
	if m.TotalInputRegister() != 90 {
		t.Errorf("Expected input register total 90 got %d", m.TotalInputRegister())
	}
	if m.mainRegister[C10] != 0 {
		t.Errorf("Expected main register coin 10 count 0, got %d", m.mainRegister[C10])
	}
	if m.mainRegister[C100] != 2 {
		t.Errorf("Expected main register coin 100 count 2, got %d", m.mainRegister[C100])
	}
	if m.inventories[1].Stock != 98 {
		t.Errorf("Expected inventories 1 stocks reduced to 98, got %d", m.inventories[1].Stock)
	}
	if m.outlet[0].Name != "Item 2" {
		t.Errorf("Expected outlet name Item 2, got %v", m.outlet[0].Name)
	}
}

func TestMachineBuySuccessfulOverpay(t *testing.T) {
	m := &Machine{
		mainRegister:  map[Currency]int{C10: 9},
		inputRegister: []Currency{C100, C100, C100},
		inventories: []Inventory{
			Inventory{
				Item{
					Name:  "Item 1",
					Price: 100,
				},
				99,
			},
			Inventory{
				Item{
					Name:  "Item 2",
					Price: 110,
				},
				99,
			},
		},
	}

	err := m.Buy(1)
	if err != nil {
		t.Errorf("Expected error nil, got %v", err.Error())
	}

	if len(m.inputRegister) != 10 {
		t.Errorf("Expected input register contain 10 coins, got %d", len(m.inputRegister))
	}
	if m.TotalInputRegister() != 190 {
		t.Errorf("Expected input register total 190 got %d", m.TotalInputRegister())
	}
	if m.mainRegister[C10] != 0 {
		t.Errorf("Expected main register coin 10 count 0, got %d", m.mainRegister[C10])
	}
	if m.mainRegister[C100] != 2 {
		t.Errorf("Expected main register coin 100 count 2, got %d", m.mainRegister[C100])
	}
	if m.inventories[1].Stock != 98 {
		t.Errorf("Expected inventories 1 stocks reduced to 98, got %d", m.inventories[1].Stock)
	}
	if m.outlet[0].Name != "Item 2" {
		t.Errorf("Expected outlet name Item 2, got %v", m.outlet[0].Name)
	}
}

func TestMachineBuySuccessfulOverpayWith500Coin(t *testing.T) {
	m := &Machine{
		mainRegister:  map[Currency]int{C10: 9, C100: 4},
		inputRegister: []Currency{C100, C500},
		inventories: []Inventory{
			Inventory{
				Item{
					Name:  "Item 1",
					Price: 100,
				},
				99,
			},
			Inventory{
				Item{
					Name:  "Item 2",
					Price: 110,
				},
				99,
			},
		},
	}

	err := m.Buy(1)
	if err != nil {
		t.Errorf("Expected error nil, got %v", err.Error())
	}

	// 4 X 100 + 9 + 10
	if len(m.inputRegister) != 13 {
		t.Errorf("Expected input register contain 13 coins, got %d", len(m.inputRegister))
	}
	if m.TotalInputRegister() != 490 {
		t.Errorf("Expected input register total 490 got %d", m.TotalInputRegister())
	}
	if m.mainRegister[C10] != 0 {
		t.Errorf("Expected main register coin 10 count 0, got %d", m.mainRegister[C10])
	}
	if m.mainRegister[C100] != 1 {
		t.Errorf("Expected main register coin 100 count 1, got %d", m.mainRegister[C100])
	}
	if m.mainRegister[C500] != 1 {
		t.Errorf("Expected main register coin 500 count 1, got %d", m.mainRegister[C500])
	}
	if m.inventories[1].Stock != 98 {
		t.Errorf("Expected inventories 1 stocks reduced to 98, got %d", m.inventories[1].Stock)
	}
	if m.outlet[0].Name != "Item 2" {
		t.Errorf("Expected outlet name Item 2, got %v", m.outlet[0].Name)
	}
}

func TestMachineDoubleBuyOverpayWithDouble500Coins(t *testing.T) {
	m := &Machine{
		mainRegister:  map[Currency]int{C10: 9, C100: 4},
		inputRegister: []Currency{C500, C500},
		inventories: []Inventory{
			Inventory{
				Item{
					Name:  "Item 1",
					Price: 100,
				},
				99,
			},
			Inventory{
				Item{
					Name:  "Item 2",
					Price: 110,
				},
				99,
			},
		},
	}

	err := m.Buy(1)
	if err != nil {
		t.Errorf("Expected error nil, got %v", err.Error())
	}
	// 1 x 500 + 3 X 100 + 9 X 10
	if len(m.inputRegister) != 13 {
		t.Errorf("Expected input register contain 13 coins, got %d", len(m.inputRegister))
	}
	if m.TotalInputRegister() != 890 {
		t.Errorf("Expected input register total 490 got %d", m.TotalInputRegister())
	}
	if m.mainRegister[C10] != 0 {
		t.Errorf("Expected main register coin 10 count 0, got %d", m.mainRegister[C10])
	}
	if m.mainRegister[C100] != 1 {
		t.Errorf("Expected main register coin 100 count 1, got %d", m.mainRegister[C100])
	}
	if m.mainRegister[C500] != 1 {
		t.Errorf("Expected main register coin 500 count 1, got %d", m.mainRegister[C500])
	}
	if m.inventories[1].Stock != 98 {
		t.Errorf("Expected inventories 1 stocks reduced to 98, got %d", m.inventories[1].Stock)
	}

	// 2nd buy
	err = m.Buy(1)
	if err != nil {
		t.Errorf("Expected error nil, got %v", err.Error())
	}
	// 1 x 500 + 2 X 100 + 8 X 10
	if len(m.inputRegister) != 11 {
		t.Errorf("Expected input register contain 13 coins, got %d", len(m.inputRegister))
	}
	if m.TotalInputRegister() != 780 {
		t.Errorf("Expected input register total 890 got %d", m.TotalInputRegister())
	}
	if m.inventories[1].Stock != 97 {
		t.Errorf("Expected inventories 1 stocks reduced to 97, got %d", m.inventories[1].Stock)
	}
	if m.mainRegister[C10] != 1 {
		t.Errorf("Expected main register coin 10 count 0, got %d", m.mainRegister[C10])
	}
	if m.mainRegister[C100] != 2 {
		t.Errorf("Expected main register coin 100 count 2, got %d", m.mainRegister[C100])
	}
	if m.mainRegister[C500] != 1 {
		t.Errorf("Expected main register coin 500 count 1, got %d", m.mainRegister[C500])
	}
	if len(m.outlet) != 2 {
		t.Errorf("Expected outlet len 2, got %d", len(m.outlet))
	}
	if m.outlet[0].Name != "Item 2" {
		t.Errorf("Expected outlet name Item 2, got %v", m.outlet[0].Name)
	}
	if m.outlet[1].Name != "Item 2" {
		t.Errorf("Expected outlet name Item 2, got %v", m.outlet[1].Name)
	}
}

func TestMachineDoubleBuyOverpayNotEnoughChangeWith500Coin(t *testing.T) {
	m := &Machine{
		mainRegister:  map[Currency]int{C10: 9, C100: 4},
		inputRegister: []Currency{C500},
		inventories: []Inventory{
			Inventory{
				Item{
					Name:  "Item 1",
					Price: 50,
				},
				99,
			},
			Inventory{
				Item{
					Name:  "Item 2",
					Price: 110,
				},
				99,
			},
		},
	}

	err := m.Buy(1)
	if err != nil {
		t.Errorf("Expected error nil, got %v", err.Error())
	}
	// 3 X 100 + 9 X 10
	if len(m.inputRegister) != 12 {
		t.Errorf("Expected input register contain 12 coins, got %d", len(m.inputRegister))
	}
	if m.TotalInputRegister() != 390 {
		t.Errorf("Expected input register total 390 got %d", m.TotalInputRegister())
	}
	if m.mainRegister[C10] != 0 {
		t.Errorf("Expected main register coin 10 count 0, got %d", m.mainRegister[C10])
	}
	if m.mainRegister[C100] != 1 {
		t.Errorf("Expected main register coin 100 count 1, got %d", m.mainRegister[C100])
	}
	if m.mainRegister[C500] != 1 {
		t.Errorf("Expected main register coin 500 count 1, got %d", m.mainRegister[C500])
	}
	if m.inventories[1].Stock != 98 {
		t.Errorf("Expected inventories 1 stocks reduced to 98, got %d", m.inventories[1].Stock)
	}

	// 2nd buy
	err = m.Buy(0)
	if err != nil {
		t.Errorf("Expected error nil, got %v", err.Error())
	}
	// 3 X 100 + 4 X 10
	if len(m.inputRegister) != 7 {
		t.Errorf("Expected input register contain 7 coins, got %d", len(m.inputRegister))
	}
	if m.TotalInputRegister() != 340 {
		t.Errorf("Expected input register total 340 got %d", m.TotalInputRegister())
	}
	if m.inventories[0].Stock != 98 {
		t.Errorf("Expected inventories 0 stocks reduced to 98, got %d", m.inventories[0].Stock)
	}
	if m.mainRegister[C10] != 5 {
		t.Errorf("Expected main register coin 10 count 1, got %d", m.mainRegister[C10])
	}
	if m.mainRegister[C100] != 1 {
		t.Errorf("Expected main register coin 100 count 2, got %d", m.mainRegister[C100])
	}
	if m.mainRegister[C500] != 1 {
		t.Errorf("Expected main register coin 500 count 1, got %d", m.mainRegister[C500])
	}
	if len(m.outlet) != 2 {
		t.Errorf("Expected outlet len 2, got %d", len(m.outlet))
	}
	if m.outlet[0].Name != "Item 2" {
		t.Errorf("Expected outlet name Item 2, got %v", m.outlet[0].Name)
	}
	if m.outlet[1].Name != "Item 1" {
		t.Errorf("Expected outlet name Item 1, got %v", m.outlet[1].Name)
	}
}

func TestMachineTotalInputRegister(t *testing.T) {
	m := &Machine{
		inputRegister: []Currency{
			C10,
			C50,
			C100,
			C100,
			C50,
			C500,
			C100,
		},
	}

	actual := m.TotalInputRegister()
	if actual != 910 {
		t.Errorf("Expected %d, got %d", 910, actual)
	}
}

func createEmptyMachine() *Machine {
	return New(map[Currency]int{}, []Inventory{})
}
