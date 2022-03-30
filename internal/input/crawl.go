package input

import (
	"fmt"
	"time"
)

func (i *FileInput) checkDisk() {
	fmt.Printf("input [%s] started working\n", i.Name)
	i.working = true
	go i.crawl()
}

func (i *FileInput) crawl() {
	time.Sleep(3 * time.Second)
	fmt.Printf("input [%s] finished working\n", i.Name)

	// in case the file input component got paused during its work
	// we do not wish to snooze, but just stay paused
	if i.state != Paused {
		i.setState(Paused)
		go i.snooze()
		// if the component is paused, we will instantly cancel the wait for the snooze func
	} else {
		i.wg.Done()
	}

	i.working = false
	i.wg.Done()
}
