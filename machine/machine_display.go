package machine

import (
	"fmt"
	"strings"
)

const MACHINE_DISPLAY_TMPL = `
[Input amount]			{INPUT}
[Change]				100 JPY			{CH_100}
						10 JPY			{CH_10}
[Return gate]			{RETURN}
[Items for sale]
{INVENTORIES}
[Outlet]				{OUTLET}
`

func (m *Machine) Display() string {
	totalInput := m.TotalInputRegister()
	input := fmt.Sprintf("%d %s", totalInput, CUR_SYMBOL)

	ch100, ch10 := m.displayChangeStatus()
	returnStr, outletStr := m.displayReturnAndOutlet()
	inventories := m.displayInventories(totalInput)

	r := strings.NewReplacer(
		"{INPUT}", input,
		"{CH_100}", ch100,
		"{CH_10}", ch10,
		"{RETURN}", returnStr,
		"{INVENTORIES}", inventories,
		"{OUTLET}", outletStr,
	)

	return strings.TrimSpace(r.Replace(MACHINE_DISPLAY_TMPL))
}

func (m *Machine) displayChangeStatus() (string, string) {
	ch100 := "No Change"
	if m.mainRegister[C100] >= 4 {
		ch100 = "Change"
	}
	ch10 := "No Change"
	if m.mainRegister[C10] >= 9 {
		ch10 = "Change"
	}

	return ch100, ch10
}

func (m *Machine) displayReturnAndOutlet() (string, string) {
	returnStr := "Empty"
	for i, v := range m.returnRegister {
		if i != 0 {
			returnStr += ", "
		} else {
			returnStr = ""
		}
		returnStr += v.Str()
	}
	outletStr := "Empty"
	for i, v := range m.outlet {
		if i != 0 {
			outletStr += ", "
		} else {
			outletStr = ""
		}
		outletStr += v.Name
	}

	return returnStr, outletStr
}

func (m *Machine) displayInventories(totalInput int) string {
	inventories := ""
	for i, v := range m.inventories {
		status := ""
		if v.Stock == 0 {
			status = "Sold out"
		} else {
			if totalInput >= v.Price {
				status = "Available for purchase"
			}
		}
		inventories += fmt.Sprintf(
			"%d. %s\t\t%d %s",
			(i + 1),
			v.Name,
			v.Price,
			CUR_SYMBOL,
		)

		if status != "" {
			inventories += fmt.Sprintf("\t\t\t%s", status)
		}

		if i != len(m.inventories)-1 {
			inventories += "\n"
		}
	}

	return inventories
}
