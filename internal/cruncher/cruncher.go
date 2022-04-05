package cruncher

import (
	"fmt"
	"sync"
)

type Data struct {
	fileName string
	fileData *[]byte
}

type Cruncher struct {
	Name     string
	arity    int
	dataChan chan *Data
	doneChan chan struct{}
	wg       sync.WaitGroup
}

func NewCruncher(name string, arity int) *Cruncher {
	return &Cruncher{
		Name:     name,
		arity:    arity,
		dataChan: make(chan *Data),
		doneChan: make(chan struct{}),
	}
}

func NewCruncherData(fileName string, fileData *[]byte) *Data {
	return &Data{
		fileName: fileName,
		fileData: fileData,
	}
}

func (c *Cruncher) Run() {
	c.wg.Add(1)
	fmt.Printf("cruncher [%s] running\n", c.Name)
	for {
		select {
		case data := <-c.dataChan:
			fmt.Printf("got data [file:%s - len:%dMB]\n", data.fileName, len(*data.fileData)/1000000)
		case <-c.doneChan:
			c.wg.Done()
			fmt.Printf("cruncher [%s] stopped\n", c.Name)
			break
		}
	}
}

func (c *Cruncher) GetDataChan() chan<- *Data {
	return c.dataChan
}

func (c *Cruncher) ShutDownGracefully(wg *sync.WaitGroup) {
	defer wg.Done()
	c.doneChan <- struct{}{}
	c.wg.Wait()
}
