package cruncher

import (
	"fmt"
	"runtime"
	"sort"
	"strings"
	"sync"
)

const counterPoolLimit = 100

type job struct {
	fileName   string
	fileData   *[]byte
	start, end int
	arity      int
	// outputs []output.Output
}

type Counter struct {
	poolChan chan struct{}
	jobChan  chan *job
	doneChan chan struct{}

	m  []map[string]int
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
	fmt.Printf("counting file [%s] chunk range [%d-%d] bytes\n", jb.fileName, jb.start, jb.end)

	m := make(map[string]int)
	fileDataLen := len(*jb.fileData)
loop:
	for i := 0; i < fileDataLen; i++ {
		if (*jb.fileData)[i] == ' ' || i == 0 {
			spaces := 0
			for j := i + 1; ; j++ {
				if j == fileDataLen || (*jb.fileData)[j] == ' ' {
					spaces++
					if spaces == jb.arity {
						start := i
						if i != 0 {
							start++
						}
						key := strings.Split(string((*jb.fileData)[start:j]), " ")
						sort.Strings(key)
						m[strings.Join(key, "-")]++
						break
					}
					if j == fileDataLen {
						break loop
					}
				}
			}
		}
	}
	fmt.Printf("counted file [%s] chunk range[%d-%d] bytes\n", jb.fileName, jb.start, jb.end)
	c.m = append(c.m, m)
	fmt.Printf("collector start")
	runtime.GC()
	fmt.Printf("collector end")
	//send output data
	c.poolChan <- struct{}{}
}
