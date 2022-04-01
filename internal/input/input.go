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
	Name      string
	state     State
	stateLock sync.Mutex
	working   bool
	sleeping  bool
	sleepDur  time.Duration
	runChan   chan State
	sleepChan chan struct{}
	doneChan  chan struct{}
	wg        sync.WaitGroup
}

func NewFileInput(name string, sleepDur time.Duration) *FileInput {
	return &FileInput{
		Name:      name,
		state:     Paused,
		sleepDur:  sleepDur,
		runChan:   make(chan State, 1),
		sleepChan: make(chan struct{}, 1),
		doneChan:  make(chan struct{}, 1),
	}
}

func (i *FileInput) Run() {
	for {
		select {
		case state := <-i.runChan:
			i.setState(state)
			// this could be set inside the switch cases so it doesn't shut down the snooze method everytime
			i.stopSnooze() // this informs the sleep routine to shut down. any instruction should reset it.
			switch i.state {
			case Started:
				fmt.Printf("input [%s] now working\n", i.Name)
			case Paused:
				fmt.Printf("input [%s] now paused\n", i.Name)
			case Stopped:
				fmt.Printf("input [%s] stopping..\n", i.Name)
				i.wg.Wait()
				// wait group or channel to await the job to finish
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

			// the idea for now is to init a job here, and when that job finishes
			// it will call the snooze function to pause the file input component for the configured amount
			// of course, if a state instruction comes from the cli, it will interrupt the snooze
			// and shut it down
			if i.working {
				break
			}
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

func (i *FileInput) setState(state State) {
	i.stateLock.Lock()
	defer i.stateLock.Unlock()
	i.state = state
}

func (i *FileInput) GetState() State {
	return i.state
}

func (i *FileInput) SendState(state State) {
	// if the channel is already full or the component is stopped, return so no deadlock appears
	if len(i.runChan) > 0 || i.state == Stopped {
		return
	}
	i.runChan <- state
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
