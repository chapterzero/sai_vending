package machine

import (
	"testing"
)

func TestNewCurrencyFromString(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected Currency
	}{
		{
			name:     "10 Coin",
			input:    "10",
			expected: C10,
		},
		{
			name:     "50 Coin",
			input:    "50",
			expected: C50,
		},
		{
			name:     "100 Coin",
			input:    "100",
			expected: C100,
		},
		{
			name:     "500 Coin",
			input:    "500",
			expected: C500,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := NewCurrencyFromString(tc.input)
			if actual != tc.expected {
				t.Errorf("Expected %d, got %d", tc.expected, actual)
			}
			if err != nil {
				t.Errorf("Expected error nil but got %v", err)
			}
		})
	}
}

func TestNewCurrencyFromStringInvalidInput(t *testing.T) {
	_, err := NewCurrencyFromString("20")
	if err == nil {
		t.Errorf("Expected error not nil, got nil")
		return
	}
	if err.Error() != "20 is not a valid coin" {
		t.Errorf("Expected error message '20 is not a valid coin' got '%s'", err.Error())
	}
}
