package output

import (
	"fmt"
	"sync"
	"time"
)

const poolLimit = 20

type jobUnion struct {
	fileName string
	chunk    map[string]int
	wg       *sync.WaitGroup
}

type jobSum struct {
	sumName   string
	fileNames []string
}

type field struct {
	bowFreq int
	mu      sync.Mutex
}

type file struct {
	bow map[string]*field
	wg  *sync.WaitGroup
}

type Info struct {
	FileName string
	Wg       *sync.WaitGroup
}

type Cache struct {
	sortProgressLimit int
	m                 map[string]*file
	poolChan          chan struct{}
	infoChan          chan *Info
	jobUnionChan      chan *jobUnion
	wg                sync.WaitGroup
}

func NewCache(sortProgressLimit int) *Cache {
	return &Cache{
		sortProgressLimit: sortProgressLimit,
		m:                 make(map[string]*file),
		poolChan:          createPool(),
		infoChan:          make(chan *Info),
		jobUnionChan:      make(chan *jobUnion),
	}
}

func (c *Cache) Run() {
	fmt.Printf("output cache running\n")
	for {
		select {
		case i := <-c.infoChan:
			_, exists := c.m[i.FileName]
			if exists {
				// wait for goroutines working on the same map field to avoid nil pointers
				c.m[i.FileName].wg.Wait()
			}
			c.m[i.FileName] = &file{
				bow: make(map[string]*field),
				wg:  i.Wg,
			}
			fmt.Printf("cache output noted of file [%s]\n", i.FileName)
		case ju := <-c.jobUnionChan:
			<-c.poolChan
			go c.unite(ju)
		}
	}
}

func (c *Cache) unite(ju *jobUnion) {
	time.Sleep(2 * time.Second)
	fmt.Printf("united chunk to file [%s]\n", ju.fileName)
	c.poolChan <- struct{}{}
}

func (c *Cache) GetInfoChan() chan<- *Info {
	return c.infoChan
}

func (c *Cache) GetJobUnionChan() chan<- *jobUnion {
	return c.jobUnionChan
}

func createPool() chan struct{} {
	c := make(chan struct{}, poolLimit)
	for i := 0; i < poolLimit; i++ {
		c <- struct{}{}
	}
	return c
}
