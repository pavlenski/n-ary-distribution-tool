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
	Name        string
	arity       int
	dataLimit   int
	dataChan    chan *Data
	doneChan    chan struct{}
	counterChan chan<- *job
	wg          sync.WaitGroup
}

func NewCruncher(name string, arity int, dataLimit int, counterChan chan<- *job) *Cruncher {
	return &Cruncher{
		Name:        name,
		arity:       arity,
		dataLimit:   dataLimit,
		dataChan:    make(chan *Data),
		doneChan:    make(chan struct{}),
		counterChan: counterChan,
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
			c.createJobs(data)
			fmt.Printf("got data [file:%s - len:%dMB]\n", data.fileName, len(*data.fileData)/1000000)
		case <-c.doneChan:
			c.wg.Done()
			break
		}
	}
}

type chunk struct {
	start, end int
}

func (c *Cruncher) createJobs(d *Data) {
	var chunks []chunk
	fileDataLen := len(*d.fileData)
	start, end := 0, 0
	cornerWords := 0
	for i := c.dataLimit; i < fileDataLen; {
		if (*d.fileData)[i] == ' ' {
			end = i
			if c.arity > 1 {
				for j := i + 1; ; j++ {
					if (*d.fileData)[j] == ' ' {
						cornerWords++
						if cornerWords == c.arity-1 {
							end = j
							break
						}
					}
				}
			}
			chunks = append(chunks, chunk{start: start, end: end})
			start = i + 1
			i += c.dataLimit
		} else {
			i++
		}
	}
	if start != fileDataLen && start != fileDataLen-1 {
		chunks = append(chunks, chunk{start: start, end: fileDataLen})
	}
	for _, chunk := range chunks {
		j := &job{
			fileName: d.fileName,
			fileData: d.fileData,
			start:    chunk.start,
			end:      chunk.end,
			arity:    c.arity,
		}
		c.counterChan <- j
	}
}

func (c *Cruncher) GetDataChan() chan<- *Data {
	return c.dataChan
}

func (c *Cruncher) ShutDownGracefully(wg *sync.WaitGroup) {
	defer wg.Done()
	c.doneChan <- struct{}{}
	c.wg.Wait()
	fmt.Printf("cruncher [%s] stopped\n", c.Name)
}
