package input

import (
	"fmt"
	"github.com/pavlenski/n-ary-distribution-tool/internal/cruncher"
	"runtime"
	"strings"
	"sync"
	"time"
)

type State = int

const (
	Started State = 1
	Paused  State = 2
	Stopped State = 0
)

type FileInput struct {
	Name             string
	disc             string
	directories      []string
	recentlyModified map[string]time.Time

	crunchers map[string]chan<- *cruncher.Data

	state     State
	stateLock sync.Mutex
	sleeping  bool
	sleepDur  time.Duration

	runChan   chan State
	sleepChan chan struct{}
	jobChan   chan *job
	wg        sync.WaitGroup
}

func NewFileInput(name, disc string, pool chan *job, sleepDur time.Duration) *FileInput {
	return &FileInput{
		Name:             name,
		disc:             disc,
		recentlyModified: make(map[string]time.Time),
		crunchers:        make(map[string]chan<- *cruncher.Data),
		state:            Paused,
		sleepDur:         sleepDur,
		runChan:          make(chan State, 1),
		sleepChan:        make(chan struct{}, 1),
		jobChan:          pool,
		wg:               sync.WaitGroup{},
	}
}

func (i *FileInput) Run() {
	fmt.Printf("running input [%s] on disk [%s]\n", i.Name, i.disc)
	for {
		select {
		case state := <-i.runChan:
			i.setState(state)
			i.stopSnooze()
			switch i.state {
			case Started:
				fmt.Printf("input [%s] now working\n", i.Name)
			case Paused:
				fmt.Printf("input [%s] now paused\n", i.Name)
			case Stopped:
				fmt.Printf("input [%s] stopping..\n", i.Name)
				return
			default:
				fmt.Printf("invalid state sent for input [%s]\n", i.Name)
			}
		default:
			runtime.Gosched()
			if i.state == Paused || i.state == Stopped {
				break
			}
			// adding one for the snooze func
			i.wg.Add(1)
			i.crawl()
		}
	}
}

// sleep will 'wake up' or unpause the file input component after configured time...
// though, if it gets informed (through the sleep channel) that it should cancel out, it does so.
func (i *FileInput) snooze() {
	i.sleeping = true
	select {
	case <-i.sleepChan:
		break
	case <-time.Tick(i.sleepDur):
		i.setState(Started)
		break
	}
	i.sleeping = false
	i.wg.Done()
}

func (i *FileInput) SendState(state State) {
	// if the channel is already full or the component is stopped, return so no deadlock appears
	if len(i.runChan) > 0 || i.state == Stopped {
		return
	}
	i.runChan <- state
}

func (i *FileInput) AddDir(dirPath string) {
	fullPath := i.disc + dirPath
	for _, dir := range i.directories {
		if fullPath == dir {
			fmt.Printf("dir [%s] already exists in input [%s]\n", fullPath, i.Name)
			return
		}
	}
	i.directories = append(i.directories, i.disc+dirPath)
	fmt.Printf("added dir [%s] to input [%s]\n", fullPath, i.Name)
}

func (i *FileInput) RemoveDir(dirPath string) {
	found := false
	fullPath := i.disc + dirPath
	dirLen := len(i.directories)
	for index, dir := range i.directories {
		if fullPath == dir {
			i.directories[index] = i.directories[dirLen-1]
			i.directories = i.directories[:dirLen-1]
			found = true
			i.clearRecentlyModified(fullPath)
			break
		}
	}
	if !found {
		fmt.Printf("directory [%s] in input [%s] does not exist\n", fullPath, i.Name)
	}
}

func (i *FileInput) LinkCruncher(c *cruncher.Cruncher) {
	if _, exists := i.crunchers[c.Name]; exists {
		fmt.Printf("curncher [%s] is already linked\n", c.Name)
		return
	}
	i.crunchers[c.Name] = c.GetDataChan()
	fmt.Printf("linking cruncher [%s] to input [%s]\n", c.Name, i.Name)
}

func (i *FileInput) clearRecentlyModified(dirPath string) {
	for filePath := range i.recentlyModified {
		if strings.HasPrefix(filePath, dirPath) {
			i.recentlyModified[filePath] = time.Time{}
		}
	}
}

func (i *FileInput) setState(state State) {
	i.stateLock.Lock()
	defer i.stateLock.Unlock()
	i.state = state
}

func (i *FileInput) stopSnooze() {
	// if the component is asleep or stopped the snooze func should be shut down
	if i.sleeping || (i.state == Stopped) {
		i.sleepChan <- struct{}{}
	}
}

func (i *FileInput) ShutDownGracefully(cliWg *sync.WaitGroup) {
	defer cliWg.Done()
	i.SendState(Stopped)
	i.wg.Wait()
	fmt.Printf("input [%s] stopped\n", i.Name)
}
