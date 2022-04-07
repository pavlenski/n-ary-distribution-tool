package app

import (
	"bufio"
	"fmt"
	"github.com/pavlenski/n-ary-distribution-tool/internal/cruncher"
	"github.com/pavlenski/n-ary-distribution-tool/internal/input"
	"github.com/pavlenski/n-ary-distribution-tool/internal/output"
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
	outputComp   = "output"

	dir = "dir"
)

func (a *App) run() {
	buffer := bufio.NewReader(os.Stdin)

	o1 := output.NewOutput("main", a.outputCache.GetJobUnionChan())
	a.outputComponents[o1.Name] = o1
	go o1.Run()

	i1 := input.NewFileInput("i1", a.discs[1], a.fileLoader.GetJobChan(a.discs[1]), a.fileInputSleepTime)
	go i1.Run()
	i1.AddDir("TEMP")
	a.inputComponents[i1.Name] = i1

	c := cruncher.NewCruncher("c1", 1, a.counterDataLimit, a.cruncherCounter.GetJobChan(), a.outputCache.GetInfoChan())
	go c.Run()
	a.cruncherComponents[c.Name] = c

	c.LinkOutput(o1)

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

			for _, c := range a.cruncherComponents {
				inputWg.Add(1)
				go c.ShutDownGracefully(inputWg)
			}
			// shut down file loader (wait for the loading to finish, then exit)
			inputWg.Add(1)
			go a.fileLoader.ShutDownGracefully(inputWg)
			inputWg.Add(1)
			go a.cruncherCounter.ShutDownGracefully(inputWg)

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
