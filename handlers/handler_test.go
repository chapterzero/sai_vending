package handlers

import (
	"testing"

	"github.com/chapterzero/sai_vending/machine"
)

func TestInsertHandle(t *testing.T) {
	testCases := []struct {
		name                 string
		cmd                  []string
		expectedErrorMessage string
	}{
		{
			name:                 "Missing required argument",
			cmd:                  []string{"1"},
			expectedErrorMessage: "Command 1 (PUSH) need 2nd argument: coin, example: 1 50",
		},
		{
			name:                 "Invalid coin",
			cmd:                  []string{"1", "30"},
			expectedErrorMessage: "30 is not a valid coin",
		},
		{
			name:                 "Successful",
			cmd:                  []string{"1", "10"},
			expectedErrorMessage: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := machine.New(map[machine.Currency]int{}, []machine.Inventory{})
			h := &InsertHandler{}
			err := h.Handle(m, tc.cmd)
			if tc.expectedErrorMessage == "" {
				if err != nil {
					t.Errorf("Expected got nil error, got %s", err.Error())
				}
			} else {
				if err.Error() != tc.expectedErrorMessage {
					t.Errorf("Expected error message %s, got %s", tc.expectedErrorMessage, err.Error())
				}
			}
		})
	}
}

func TestBuyHandle(t *testing.T) {
	testCases := []struct {
		name                 string
		cmd                  []string
		expectedErrorMessage string
	}{
		{
			name:                 "Missing required argument",
			cmd:                  []string{"2"},
			expectedErrorMessage: "Command 2 (BUY) need 2nd argument: #item, example: 2 1 to buy first item",
		},
		{
			name:                 "Invalid item index",
			cmd:                  []string{"2", "a"},
			expectedErrorMessage: "strconv.Atoi: parsing \"a\": invalid syntax",
		},
		{
			name:                 "Successful",
			cmd:                  []string{"2", "1"},
			expectedErrorMessage: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := machine.New(map[machine.Currency]int{}, []machine.Inventory{
				machine.Inventory{
					machine.Item{
						Name:  "Item 1",
						Price: 10,
					},
					99,
				},
			})

			m.Insert(machine.C10)

			h := &BuyHandler{}
			err := h.Handle(m, tc.cmd)
			if tc.expectedErrorMessage != "" {
				if err.Error() != tc.expectedErrorMessage {
					t.Errorf("Expected error message '%s', got '%s'", tc.expectedErrorMessage, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected got nil error, got %s", err.Error())
				}
			}
		})
	}
}

func TestGetItemHandler(t *testing.T) {
	m := machine.New(map[machine.Currency]int{machine.C10: 9, machine.C100: 4}, []machine.Inventory{
		machine.Inventory{
			machine.Item{
				Name:  "Item 1",
				Price: 10,
			},
			99,
		},
	})

	m.Insert(machine.C100)
	m.Buy(0)
	m.Buy(0)

	h := &GetItemHandler{}
	err := h.Handle(m, []string{})

	if err != nil {
		t.Errorf("Expected got nil error, got %s", err.Error())
	}
}

func TestReturnInputHandler(t *testing.T) {
	m := machine.New(map[machine.Currency]int{machine.C10: 9, machine.C100: 4}, []machine.Inventory{
		machine.Inventory{
			machine.Item{
				Name:  "Item 1",
				Price: 10,
			},
			99,
		},
	})

	m.Insert(machine.C100)
	m.Buy(0)

	h := &ReturnInputHandler{}
	err := h.Handle(m, []string{})

	if err != nil {
		t.Errorf("Expected got nil error, got %s", err.Error())
	}
}

func TestGetReturnHandler(t *testing.T) {
	m := machine.New(map[machine.Currency]int{machine.C10: 9, machine.C100: 4}, []machine.Inventory{
		machine.Inventory{
			machine.Item{
				Name:  "Item 1",
				Price: 10,
			},
			99,
		},
	})

	m.Insert(machine.C100)
	m.Buy(0)
	m.ReturnInput()

	h := &GetReturnHandler{}
	err := h.Handle(m, []string{})
	if err != nil {
		t.Errorf("Expected got nil error, got %s", err.Error())
	}
}
