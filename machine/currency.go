package machine

import (
	"fmt"
)

type Currency int

const CUR_SYMBOL = "JPY"

const (
	C10  Currency = 10
	C50  Currency = 50
	C100 Currency = 100
	C500 Currency = 500
)

func NewCurrencyFromString(s string) (Currency, error) {
	switch s {
	case "10":
		return C10, nil
	case "50":
		return C50, nil
	case "100":
		return C100, nil
	case "500":
		return C500, nil
	}

	return Currency(-1), fmt.Errorf("%s is not a valid coin", s)
}
