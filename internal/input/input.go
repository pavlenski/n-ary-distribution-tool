package input

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type State = int

const (
	Working State = 1
	Paused  State = 2
	Stopped State = 0
)

type FileInput struct {
	Name string

	state     State
	stateLock sync.Mutex

	sleepDur time.Duration

	runChan   <-chan State
	sleepChan chan struct{}
}

func NewFileInput(name string, sleepDur time.Duration) *FileInput {
	return &FileInput{
		Name:      name,
		state:     Working,
		sleepDur:  sleepDur,
		runChan:   make(<-chan State),
		sleepChan: make(chan struct{}),
	}
}

func (i *FileInput) Run() {
	for {
		select {
		case state := <-i.runChan:
			i.setState(state)
			i.sleepChan <- struct{}{} // this informs the sleep routine to shut down. any instruction should reset it.
			switch i.state {
			case Working:
				fmt.Printf("input [%s] now working", i.Name)
			case Paused:
				fmt.Printf("input [%s] now paused", i.Name)
			case Stopped:
				fmt.Printf("input [%s] stopping..\n", i.Name)
				// wait group or channel to await the job to finish
				fmt.Printf("input [%s] now stopped\n", i.Name)
				return
			}
		default:
			runtime.Gosched()
			if i.state == Paused {
				break
			}

			// init job

			// the idea for now is to init a job here, and when that job finishes
			// it will call the snooze function to pause the file input component for the configured amount
			// of course, if a state instruction comes from the cli, it will interrupt the snooze
			// and shut it down
		}
	}
}

// sleep will 'wake up' or unpause the file input component after configured time...
// though, if it gets informed (through the sleep channel) that it should cancel out, it does so.
func (i *FileInput) snooze() {
	select {
	case <-i.sleepChan:
		return
	case <-time.Tick(i.sleepDur):
		fmt.Printf("hey input [%s] wake up!\n", i.Name)
		i.setState(Working)
	}
}

func (i *FileInput) setState(state State) {
	i.stateLock.Lock()
	defer i.stateLock.Unlock()
	i.state = state
}

func (i *FileInput) GetState() State {
	return i.state
}
