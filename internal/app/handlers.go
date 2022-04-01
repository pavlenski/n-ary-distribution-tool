package app

import (
	"fmt"
	"github.com/pavlenski/n-ary-distribution-tool/internal/input"
)

func (a *App) handleAddCommand(component string, args []string) {
	switch component {
	case fileInput:
		a.addFileInput(args[0])
	default:
		fmt.Printf(
			"'%s' is an invalid component type.. use one of [%s | %s | %s]\n",
			component, fileInput, cruncher, output,
		)
	}
}

func (a *App) handleRemoveCommand(component string) {
	switch component {
	case fileInput:
		a.handleInputStateCommand(component, input.Stopped)
	default:
		fmt.Printf(
			"'%s' is an invalid component type.. use one of [%s | %s | %s]\n",
			component, fileInput, cruncher, output,
		)
	}
}

func (a *App) handleInputStateCommand(inputName string, state input.State) {
	i, ok := a.inputComponents[inputName]
	if !ok {
		fmt.Printf("input [%s] does not exist\n", inputName)
		return
	}
	i.SendState(state)
}

func (a *App) addFileInput(name string) {
	_, exists := a.inputComponents[name]
	if exists {
		fmt.Printf("input with name [%s] already exists", name)
	}

	i := input.NewFileInput(name, a.fileInputSleepTime)
	a.inputComponents[i.Name] = i
	fmt.Printf("adding input [%s]\n", name)
	go i.Run()
}
