package app

import (
	"fmt"
	"github.com/pavlenski/n-ary-distribution-tool/internal/input"
	"os"
	"strconv"
)

func (a *App) handleAddCommand(component string, args []string) {
	switch component {
	case fileInput:
		discIndex, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("disc index argument must be a number\n")
			return
		}
		a.addFileInput(args[0], discIndex)
	case dir:
		fi, ok := a.inputComponents[args[1]]
		if !ok {
			fmt.Printf("input [%s] does not exist.. can't add dir", args[1])
			return
		}
		fi.AddDir(args[0])
	default:
		fmt.Printf(
			"'%s' is an invalid component type.. use one of [%s | %s | %s]\n",
			component, fileInput, cruncher, output,
		)
	}
}

func (a *App) handleRemoveCommand(component string, args []string) {
	switch component {
	case fileInput:
		name := args[0]
		a.handleInputStateCommand(name, input.Stopped)
		delete(a.inputComponents, name)
	case dir:
		fi, ok := a.inputComponents[args[1]]
		if !ok {
			fmt.Printf("input [%s] does not exist.. can't add dir\n", args[1])
			return
		}
		fi.RemoveDir(args[0])
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

func (a *App) addFileInput(name string, discIndex int) {
	_, exists := a.inputComponents[name]
	if exists {
		fmt.Printf("input with name [%s] already exists\n", name)
		return
	}
	disc, ok := a.discs[discIndex]
	if !ok {
		fmt.Printf("disk of index [%d] does not exist\n", discIndex)
		return
	}
	i := input.NewFileInput(name, disc, a.fileLoader.GetJobChan(disc), a.fileInputSleepTime)
	a.inputComponents[i.Name] = i
	go i.Run()
}

func dirExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
