package app

import (
	"fmt"
	"github.com/pavlenski/n-ary-distribution-tool/internal/input"
	"time"
)

func (a *App) handleAddCommand(component string, args []string) {
	switch component {
	case fileInput:
		// TODO clean up so errors can be caught
		dur, _ := time.ParseDuration(args[1])
		_, exists := a.inputHandlers[args[0]]
		if exists {
			fmt.Printf("input with name [%s] already exists", args[1])
		}

		i := input.NewFileInput(args[0], dur)
		a.inputHandlers[i.Name] = i
		fmt.Printf("adding input [%s]\n", args[0])
		go i.Run()
	default:
		fmt.Printf(
			"'%s' is an invalid component type.. use one of [%s | %s | %s]\n",
			component, fileInput, cruncher, output,
		)
	}
}

func (a *App) handleInputStateCommand(inputName string, state input.State) {
	i, ok := a.inputHandlers[inputName]
	if !ok {
		fmt.Printf("input [%s] does not exist\n", inputName)
		return
	}
	i.SendState(state)
}
