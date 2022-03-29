package app

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	add    = "add"
	link   = "link"
	remove = "remove"
	status = "status"
	exit   = "exit"
)

const (
	input    = "input"
	cruncher = "cruncher"
	output   = "output"
)

func (a *App) run() {
	indent()
	buffer := bufio.NewReader(os.Stdin)

	for {
		line, err := buffer.ReadString('\n')
		if err != nil {
			log.Fatalln("error scanning command.. exiting.")
		}

		args := extractArgs(line)
		command := args[0]

		_ = command

		switch command {
		case add:
			a.handleAddCommand(args[1], args[2:])
		case link:
			fmt.Println("linking", args[1:])
		case remove:
			fmt.Println("removing", args[1:])
		case status:
			fmt.Println("printing status")
		case exit:
			fmt.Println("exiting")
			// send stop signals and stuff
			return
		}

	}
}

func indent() {
	fmt.Printf("cli: ")
}

func extractArgs(line string) []string {
	return strings.Split(line[:len(line)-1], " ")
}

func (a *App) handleAddCommand(component string, args []string) {
	switch component {
	case input:
		fmt.Printf("adding %s component with args %v\n", component, args)
	default:
		fmt.Printf("'%s' is an invalid component type.. use one of [%s | %s | %s]\n", component, input, cruncher, output)
	}
}
