package app

import (
	"bufio"
	"fmt"
	"github.com/pavlenski/n-ary-distribution-tool/internal/input"
	"log"
	"os"
	"runtime"
	"strings"
)

const (
	add    = "add"
	link   = "link"
	remove = "remove"
	status = "status"
	exit   = "exit"

	pause = "pause"
	start = "start"
	stop  = "stop"

	memory   = "mem"
	megabyte = 1000000
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
			a.handleRemoveCommand(args[1])
		case pause:
			a.handleInputStateCommand(args[1], input.Paused)
		case start:
			a.handleInputStateCommand(args[1], input.Started)
		case status:
			fmt.Println("printing status")
		case memory:
			mem()
		case exit:
			// send stop signals and stuff
			return
		}

	}
}

func mem() {
	ms := &runtime.MemStats{}
	runtime.ReadMemStats(ms)
	fmt.Printf(
		"HEAPALLOC: %fMB\nALLOC: %fMB\nSYS: %fMB\n",
		float64(ms.HeapAlloc)/megabyte,
		float64(ms.Alloc)/megabyte,
		float64(ms.Sys)/megabyte,
	)
}

func indent() {
	fmt.Printf("cli: ")
}

func extractArgs(line string) []string {
	return strings.Split(line[:len(line)-1], " ")
}
