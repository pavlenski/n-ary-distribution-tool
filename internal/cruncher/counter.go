package cruncher

import (
	"fmt"
	"github.com/pavlenski/n-ary-distribution-tool/internal/output"
	"sort"
	"strings"
	"sync"
)

const counterPoolLimit = 20

type job struct {
	fileName   string
	fileData   *[]byte
	start, end int
	arity      int
	outputs    []chan<- *output.Data
	wg         *sync.WaitGroup
}

type Counter struct {
	poolChan chan struct{}
	jobChan  chan *job
	doneChan chan struct{}

	wg sync.WaitGroup
}

func NewCounter() *Counter {
	return &Counter{
		poolChan: createPool(),
		jobChan:  make(chan *job, 2*counterPoolLimit),
		doneChan: make(chan struct{}),
	}
}

func (c *Counter) Run() {
	for {
		select {
		case j := <-c.jobChan:
			<-c.poolChan
			c.wg.Add(1)
			go c.countAndForward(j)
		case <-c.doneChan:
			break
		}
	}
}

func (c *Counter) GetJobChan() chan<- *job {
	return c.jobChan
}

func (c *Counter) ShutDownGracefully(wg *sync.WaitGroup) {
	defer wg.Done()
	c.doneChan <- struct{}{}
	c.wg.Wait()
	fmt.Printf("cruncher counter stopped\n")
}

func createPool() chan struct{} {
	p := make(chan struct{}, counterPoolLimit)
	for i := 0; i < counterPoolLimit; i++ {
		p <- struct{}{}
	}
	return p
}

func (c *Counter) countAndForward(jb *job) {
	defer c.wg.Done()
	//fmt.Printf("counting file [%s] chunk range [%d-%d] bytes\n", jb.fileName, jb.start, jb.end)
	//fmt.Printf("data-chunk: %s\n", (*jb.fileData)[jb.start:jb.end])

	m := make(map[string]int)
loop:
	for i := jb.start; i < jb.end; i++ {
		if (*jb.fileData)[i] == ' ' || i == jb.start {
			spaces := 0
			for j := i + 1; ; j++ {
				if j == jb.end || (*jb.fileData)[j] == ' ' {
					spaces++
					if spaces == jb.arity {
						start := i
						if i != jb.start {
							start++
						}
						key := strings.Split(string((*jb.fileData)[start:j]), " ")
						sort.Strings(key)
						m[strings.Join(key, "-")]++
						//fmt.Printf("WORD [%s] START [%d] END [%d] I [%d]\n", strings.Join(key, "-"), start, j, i)
						break
					}
					if j == jb.end {
						break loop
					}
				}
			}
		}
	}

	d := &output.Data{
		FileName: jb.fileName,
		Chunk:    m,
		FileWg:   jb.wg,
	}
	for _, o := range jb.outputs {
		o <- d
	}
	//fmt.Println("chunk", m)
	//send output data
	fmt.Printf("counted file [%s] chunk range[%d-%d] bytes\n", jb.fileName, jb.start, jb.end)
	c.poolChan <- struct{}{}
}
