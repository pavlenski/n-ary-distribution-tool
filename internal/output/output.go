package output

import (
	"fmt"
	"sync"
)

type Data struct {
	FileName string
	Chunk    map[string]int
	FileWg   *sync.WaitGroup
}

type Output struct {
	Name      string
	dataChan  chan *Data
	doneChan  chan struct{}
	unionChan chan<- *jobUnion
	wg        sync.WaitGroup
}

func NewOutput(name string, unionChan chan<- *jobUnion) *Output {
	return &Output{
		Name:      name,
		dataChan:  make(chan *Data),
		doneChan:  make(chan struct{}),
		unionChan: unionChan,
	}
}

func (o *Output) Run() {
	for {
		select {
		case data := <-o.dataChan:
			o.unionChan <- &jobUnion{
				fileName: data.FileName,
				chunk:    data.Chunk,
				wg:       data.FileWg,
			}
		case <-o.doneChan:
			break
		}
	}
}

func (o *Output) GetDataChan() chan<- *Data {
	return o.dataChan
}

func (o *Output) ShutDownGracefully(wg *sync.WaitGroup) {
	defer wg.Done()
	o.doneChan <- struct{}{}
	o.wg.Wait()
	fmt.Printf("output [%s] stopped\n", o.Name)
}
