package app

import (
	"bufio"
	"fmt"
	"github.com/pavlenski/n-ary-distribution-tool/internal/input"
	"log"
	"os"
	"strings"
	"time"
)

const (
	add    = "add"
	link   = "link"
	remove = "remove"
	status = "status"
	exit   = "exit"
)

const (
	fileInput = "input"
	cruncher  = "cruncher"
	output    = "output"
)

func (a *App) run() {
	buffer := bufio.NewReader(os.Stdin)

	for {
		indent()
		line, err := buffer.ReadString('\n')
		if err != nil {
			log.Fatalln("error scanning command.. exiting.")
		}

		args := extractArgs(line)
		command := args[0]

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
	case fileInput:
		fmt.Printf("adding %s component with args %v\n", component, args)
		dur, _ := time.ParseDuration(args[1])
		i := input.NewFileInput(args[0], dur)
		go i.Run()

	default:
		fmt.Printf(
			"'%s' is an invalid component type.. use one of [%s | %s | %s]\n",
			component, fileInput, cruncher, output,
		)
	}
}
