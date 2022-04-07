package cruncher

import (
	"fmt"
	"github.com/pavlenski/n-ary-distribution-tool/internal/output"
	"strings"
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
	outputs     map[string]chan<- *output.Data
	infoChan    chan<- *output.Info
	dataChan    chan *Data
	doneChan    chan struct{}
	counterChan chan<- *job
	wg          sync.WaitGroup
}

func NewCruncher(
	name string,
	arity int,
	dataLimit int,
	counterChan chan<- *job,
	infoChan chan<- *output.Info,
) *Cruncher {
	return &Cruncher{
		Name:        name,
		arity:       arity,
		dataLimit:   dataLimit,
		outputs:     make(map[string]chan<- *output.Data),
		infoChan:    infoChan,
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
			fmt.Printf("cruncher [%s] got data [file:%s - len:%dMB]\n", c.Name, data.fileName, len(*data.fileData)/1000000)
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
loop:
	for i := c.dataLimit; i < fileDataLen; {
		if (*d.fileData)[i] == ' ' {
			end = i
			if c.arity > 1 {
				for j := i + 1; ; j++ {
					if j == fileDataLen || (*d.fileData)[j] == ' ' {
						cornerWords++
						if cornerWords == c.arity-1 {
							end = j
							break
						}
					}
					if j == fileDataLen {
						break loop
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
	var outputs []chan<- *output.Data
	for _, o := range c.outputs {
		outputs = append(outputs, o)
	}
	// here i signal output what is worked on
	fileWg := sync.WaitGroup{}
	fileWg.Add(len(chunks))
	s := strings.Split(d.fileName, ".")
	fileName := fmt.Sprintf("%s-arity%d.%s", s[0], c.arity, s[1])
	c.infoChan <- &output.Info{
		FileName: fileName,
		Wg:       &fileWg,
	}

	for _, chunk := range chunks {
		j := &job{
			fileName: fileName,
			fileData: d.fileData,
			start:    chunk.start,
			end:      chunk.end,
			arity:    c.arity,
			outputs:  outputs,
			wg:       &fileWg,
		}
		c.counterChan <- j
	}
}

func (c *Cruncher) LinkOutput(o *output.Output) {
	c.outputs[o.Name] = o.GetDataChan()
	fmt.Printf("linking output [%s] to cruncher [%s]\n", o.Name, c.Name)
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
