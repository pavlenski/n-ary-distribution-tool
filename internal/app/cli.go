package app

import (
	"bufio"
	"fmt"
	"github.com/pavlenski/n-ary-distribution-tool/internal/input"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
)

const (
	add    = "add"
	link   = "link"
	remove = "remove"
	status = "status"
	exit   = "exit"

	pause = "pause"
	start = "start"

	memory   = "mem"
	megabyte = 1000000
)

const (
	fileInput = "input"
	cruncher  = "cruncher"
	output    = "output"

	dir = "dir"
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
			if len(args) < 4 {
				fmt.Printf("not enoguh args, try again..\n")
				continue
			}
			a.handleAddCommand(args[1], args[2:])
		case link:
			// temp
			for _, fi := range a.inputComponents {
				//fi.TempPrintDirs()
				fi.TempRecently()
			}
		case remove:
			if len(args) < 4 {
				fmt.Printf("not enoguh args, try again..\n")
				continue
			}
			a.handleRemoveCommand(args[1], args[2:])
		case pause:
			a.handleInputStateCommand(args[1], input.Paused)
		case start:
			a.handleInputStateCommand(args[1], input.Started)
		case status:
			fmt.Println("printing status")
		case memory:
			mem()
		case exit:
			inputWg := &sync.WaitGroup{}
			for _, fi := range a.inputComponents {
				inputWg.Add(1)
				go fi.ShutDown(inputWg)
			}
			inputWg.Wait()

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
