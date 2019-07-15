package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chapterzero/sai_vending/handlers"
	"github.com/chapterzero/sai_vending/machine"
)

var m *machine.Machine
var hMap map[string]handlers.Handler

func init() {
	log.Println("Initializing...")
	m = provisionMachine()
	hMap = map[string]handlers.Handler{
		"1": &handlers.InsertHandler{},
		"2": &handlers.BuyHandler{},
		"3": &handlers.GetItemHandler{},
		"4": &handlers.ReturnInputHandler{},
		"5": &handlers.GetReturnHandler{},
	}
}

func main() {
	log.Println("SAI VENDING PROGRAM v0.1 press CTRL-C to exit")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(m.Display())

	for {
		fmt.Println("--------------------------------------------------------")
		log.Println("Enter command")
		scanner.Scan()

		cmd := strings.Split(scanner.Text(), " ")
		if h, ok := hMap[cmd[0]]; ok {
			err := h.Handle(m, cmd)
			if err != nil {
				printError(err)
				continue
			}
			fmt.Println(m.Display())
			continue
		}

		printError(fmt.Errorf("Invalid commands"))
	}
}

func printError(err error) {
	log.Println("ERR:", err.Error())
}

func provisionMachine() *machine.Machine {
	return machine.New(map[machine.Currency]int{
		machine.C10:  200,
		machine.C100: 10,
	}, []machine.Inventory{
		machine.Inventory{
			machine.Item{
				Name:  "Canned Coffee",
				Price: 120,
			},
			10,
		},
		machine.Inventory{
			machine.Item{
				Name:  "Water PET bottle",
				Price: 100,
			},
			0,
		},
		machine.Inventory{
			machine.Item{
				Name:  "Sport drinks XT",
				Price: 150,
			},
			5,
		},
	})
}
