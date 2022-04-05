package app

import (
	"bufio"
	"fmt"
	"github.com/pavlenski/n-ary-distribution-tool/internal/cruncher"
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
	fileInput    = "input"
	cruncherComp = "cruncher"
	output       = "output"

	dir = "dir"
)

func (a *App) run() {
	buffer := bufio.NewReader(os.Stdin)

	c := cruncher.NewCruncher("c1", 1)
	go c.Run()
	a.cruncherComponents[c.Name] = c

	for {
		line, err := buffer.ReadString('\n')
		if err != nil {
			log.Fatalln("error scanning command.. exiting.")
		}

		args := extractArgs(line)
		command := args[0]

		switch command {
		case "temp":
		case add:
			a.handleAddCommand(args[1], args[2:])
		case link:
			a.handleLinkCommand(args[1], args[2])
		case remove:
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
			// shut down input components (prevent job creation)
			for _, fi := range a.inputComponents {
				inputWg.Add(1)
				go fi.ShutDownGracefully(inputWg)
			}
			// shut down file loader (wait for the loading to finish, then exit)
			inputWg.Add(1)
			go a.fileLoader.ShutDownGracefully(inputWg)
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

func extractArgs(line string) []string {
	return strings.Split(line[:len(line)-1], " ")
}
