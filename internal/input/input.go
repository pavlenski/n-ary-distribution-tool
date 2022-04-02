package input

import (
	"fmt"
	"runtime"
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

	state     State
	stateLock sync.Mutex
	working   bool
	sleeping  bool
	sleepDur  time.Duration

	runChan   chan State
	sleepChan chan struct{}
	doneChan  chan struct{}
	// this ought to be changed from a struct{}
	poolChan chan struct{}
	wg       sync.WaitGroup
}

func NewFileInput(name, disc string, pool chan struct{}, sleepDur time.Duration) *FileInput {
	return &FileInput{
		Name:             name,
		disc:             disc,
		recentlyModified: make(map[string]time.Time),
		state:            Paused,
		sleepDur:         sleepDur,
		runChan:          make(chan State, 1),
		sleepChan:        make(chan struct{}, 1),
		doneChan:         make(chan struct{}, 1),
		poolChan:         pool,
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
				i.wg.Wait()
				fmt.Printf("input [%s] now stopped\n", i.Name)
				i.doneChan <- struct{}{}
				return
			default:
				fmt.Printf("invalid state sent for input [%s]\n", i.Name)
			}
		default:
			runtime.Gosched()
			if i.state == Paused || i.state == Stopped {
				break
			}
			if i.working {
				break
			}

			// the idea for now is to init a job here, and when that job finishes
			// it will call the snooze function to pause the file input component for the configured amount
			// of course, if a state instruction comes from the cli, it will interrupt the snooze
			// and shut it down

			// adding one for the crawling and one for the snooze func
			i.wg.Add(2)
			i.checkDisk()
		}
	}
}

// sleep will 'wake up' or unpause the file input component after configured time...
// though, if it gets informed (through the sleep channel) that it should cancel out, it does so.
func (i *FileInput) snooze() {
	i.sleeping = true
	select {
	case <-i.sleepChan:
		fmt.Printf("stopping input [%s] snooze func\n", i.Name)
		break
	case <-time.Tick(i.sleepDur):
		fmt.Printf("hey input [%s] wake up!\n", i.Name)
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
		}
	}
	if !found {
		fmt.Printf("directory [%s] in input [%s] does not exist\n", fullPath, i.Name)
	}
}

func (i *FileInput) TempPrintDirs() {
	fmt.Println(i.directories)
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

func (i *FileInput) ShutDown(wg *sync.WaitGroup) {
	defer wg.Done()
	i.SendState(Stopped)
	<-i.doneChan
}

func (i *FileInput) TempRecently() {
	fmt.Println(i.recentlyModified)
}
