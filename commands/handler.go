package commands

import (
	"fmt"
	"strconv"

	"github.com/chapterzero/sai_vending/machine"
)

type Handler interface {
	Handle(m *machine.Machine, cmd []string) error
}

type InsertHandler struct{}

func (h *InsertHandler) Handle(m *machine.Machine, cmd []string) error {
	if len(cmd) < 2 {
		return fmt.Errorf("Command 1 (PUSH) need 2nd argument: coin, example: 1 50")
	}

	c, err := machine.NewCurrencyFromString(cmd[1])
	if err != nil {
		return err
	}

	return m.Insert(c)
}

type BuyHandler struct{}

func (h *BuyHandler) Handle(m *machine.Machine, cmd []string) error {
	if len(cmd) < 2 {
		return fmt.Errorf("Command 2 (BUY) need 2nd argument: #item, example: 2 1 to buy first item")
	}

	idx, err := strconv.Atoi(cmd[1])
	if err != nil {
		return err
	}

	return m.Buy(idx - 1)
}

type GetItemHandler struct{}

func (h *GetItemHandler) Handle(m *machine.Machine, cmd []string) error {
	str := ""
	for i, item := range m.GetItems() {
		if i != 0 {
			str = str + ", "
		}
		str += item.Name
	}

	if str != "" {
		fmt.Println("GOT Items:", str)
	}

	return nil
}

type ReturnInputHandler struct{}

func (h *ReturnInputHandler) Handle(m *machine.Machine, cmd []string) error {
	m.ReturnInput()
	return nil
}

type GetReturnHandler struct{}

func (h *GetReturnHandler) Handle(m *machine.Machine, cmd []string) error {
	str := ""
	for i, c := range m.GetReturn() {
		if i != 0 {
			str = str + ", "
		}
		str += c.Str()
	}
	if str != "" {
		fmt.Println("GOT Changes:", str)
	}

	return nil
}
