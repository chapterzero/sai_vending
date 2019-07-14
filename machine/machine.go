package machine

import (
	"fmt"
)

func New(provision map[Currency]int, inventories []Inventory) *Machine {
	return &Machine{
		mainRegister:   provision,
		inputRegister:  make([]Currency, 0),
		returnRegister: make([]Currency, 0),
		inventories:    inventories,
		outlet:         make([]Item, 0),
	}
}

type Machine struct {
	// the machine drawer
	// key is Currency, the value is the total. Ex:
	mainRegister map[Currency]int

	// money inserted to machine entered this register 1st
	inputRegister []Currency

	returnRegister []Currency
	inventories    []Inventory
	outlet         []Item
}

func (m *Machine) Insert(c Currency) error {
	// check if machine able to return
	if (c != C10 && m.mainRegister[C10] < 9) ||
		(c == C500 && m.mainRegister[C100] < 4) {
		m.returnRegister = append(m.returnRegister, c)
		return fmt.Errorf("Unable to return change")
	}

	m.inputRegister = append(m.inputRegister, c)
	return nil
}

// index start from zero
// return error to check if the buy successful / not
// (nil error for successful buy)
func (m *Machine) Buy(i int) error {
	err := m.isAllowToBuy(i)
	if err != nil {
		return err
	}

	// for rollback purpose,
	// transaction operation is not done on real register
	mR, iR := m.createRegisterCopy()

	// begin transaction
	taken := 0
	takenIdx := 0
	for _, v := range iR {
		taken += int(v)
		takenIdx++
		mR[v]++
		if taken >= m.inventories[i].Price {
			break
		}
	}
	// deduct input, calculate change
	iR = iR[takenIdx:]
	mR, iR, err = calculateChange(mR, iR, taken, m.inventories[i].Price)
	if err != nil {
		return err
	}

	// commiting transaction and
	// return the changes to input register to allow multiple buy
	// stock deduction & disperse
	m.mainRegister, m.inputRegister = mR, iR
	m.inventories[i].Stock--
	m.outlet = append(m.outlet, m.inventories[i].Item)

	return nil
}

func (m *Machine) GetItems() []Item {
	defer func() {
		m.outlet = []Item{}
	}()

	return m.outlet
}

func (m *Machine) ReturnInput() {
	m.returnRegister = append(m.returnRegister, m.inputRegister...)
	m.inputRegister = []Currency{}
}

func (m *Machine) GetReturn() []Currency {
	defer func() {
		m.returnRegister = []Currency{}
	}()
	return m.returnRegister
}

func (m *Machine) TotalInputRegister() int {
	ttl := 0
	for _, v := range m.inputRegister {
		ttl += int(v)
	}
	return ttl
}

func (m *Machine) Display() string {
}

func (m *Machine) createRegisterCopy() (map[Currency]int, []Currency) {
	mainRegister := make(map[Currency]int, 0)
	inputRegister := make([]Currency, len(m.inputRegister))

	copy(inputRegister, m.inputRegister)
	for k, v := range m.mainRegister {
		mainRegister[k] = v
	}

	return mainRegister, inputRegister
}

func (m *Machine) isAllowToBuy(i int) error {
	if i < 0 || i >= len(m.inventories) {
		return fmt.Errorf("Invalid inventory, please enter number from (1 to %d)", len(m.inventories))
	}

	ttlInput := m.TotalInputRegister()
	if ttlInput < m.inventories[i].Price {
		return fmt.Errorf("Inserted money not enough to buy this item")
	}

	return nil
}

func calculateChange(mR map[Currency]int, iR []Currency, taken, itemPrice int) (map[Currency]int, []Currency, error) {
	change := taken - itemPrice
	for {
		if change == 0 {
			break
		}

		//prioritize C100 before C10
		if change >= 100 && mR[C100] > 0 {
			change -= int(C100)
			mR[C100]--
			iR = append([]Currency{C100}, iR...)
			continue
		}

		if mR[C10] > 0 {
			change -= int(C10)
			mR[C10]--
			iR = append([]Currency{C10}, iR...)
		} else {
			// TODO: create unit test to prove if this condition might happen
			// there is still remaining change
			// but coin to return are empty
			return mR, iR, fmt.Errorf("Unable to return change")
		}
	}

	return mR, iR, nil
}
