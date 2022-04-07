package output

import (
	"fmt"
	"sort"
	"sync"
)

const poolLimit = 20

type jobUnion struct {
	fileName string
	chunk    map[string]int
	wg       *sync.WaitGroup
}

type file struct {
	bow  map[string]int
	done bool
	wg   *sync.WaitGroup
}

type Info struct {
	FileName string
	Wg       *sync.WaitGroup
}

type Cache struct {
	sortProgressLimit int
	m                 map[string]*file
	mu                *sync.RWMutex
	poolChan          chan struct{}
	infoChan          chan *Info
	jobUnionChan      chan *jobUnion
	doneChan          chan struct{}
	wg                sync.WaitGroup
}

func NewCache(sortProgressLimit int) *Cache {
	return &Cache{
		sortProgressLimit: sortProgressLimit,
		m:                 make(map[string]*file),
		mu:                &sync.RWMutex{},
		poolChan:          createPool(),
		infoChan:          make(chan *Info),
		jobUnionChan:      make(chan *jobUnion),
		doneChan:          make(chan struct{}),
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
				bow:  make(map[string]int),
				done: false,
				wg:   i.Wg,
			}
			go c.updateWhenDone(i.FileName, c.m[i.FileName].wg, &c.m[i.FileName].done)
			fmt.Printf("cache output noted of file [%s]\n", i.FileName)
		case ju := <-c.jobUnionChan:
			<-c.poolChan
			go c.unite(ju)
		case <-c.doneChan:
			break
		}
	}
}

func (c *Cache) updateWhenDone(fileName string, wg *sync.WaitGroup, done *bool) {
	wg.Wait()
	*done = true
	fmt.Printf("file [%s] finished mapping\n", fileName)
}

func (c *Cache) Sum(sumName string, fileNames []string) {
	fmt.Printf("initializing sum for [%s]\n", sumName)
	sumMap := make(map[string]int)
	for _, fn := range fileNames {
		fmt.Printf("waiting for file [%s] to finish mapping..\n", fn)
		c.m[fn].wg.Wait()
		c.mu.Lock()
		for k, v := range c.m[fn].bow {
			sumMap[k] += v
		}
		c.mu.Unlock()
	}
	c.m[sumName] = &file{
		bow:  sumMap,
		done: true,
		wg:   &sync.WaitGroup{},
	}
	fmt.Printf("sum [%s] complete\n", sumName)
}

func (c *Cache) unite(ju *jobUnion) {
	defer c.m[ju.fileName].wg.Done()
	for k, v := range ju.chunk {
		c.mu.Lock()
		c.m[ju.fileName].bow[k] += v
		c.mu.Unlock()
	}
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

func (c *Cache) Poll(fileName string) {
	if _, exists := c.m[fileName]; !exists {
		fmt.Printf("data for [%s] doesn't exist\n", fileName)
		return
	}
	if !c.m[fileName].done {
		fmt.Printf("uniting for [%s] not yet done..\n", fileName)
		return
	}
	c.sortAndPrint(fileName)
}

func (c *Cache) Take(fileName string) {
	if _, exists := c.m[fileName]; !exists {
		fmt.Printf("data for [%s] was never initiated, cannot take\n", fileName)
		return
	}
	fmt.Printf("waiting for [%s] results...\n", fileName)
	c.m[fileName].wg.Wait()
	c.sortAndPrint(fileName)
}

type result struct {
	key   string
	value int
}

func (c *Cache) sortAndPrint(fileName string) {
	c.mu.Lock()
	sorted := make([]result, 0, len(c.m[fileName].bow))
	for k, v := range c.m[fileName].bow {
		sorted = append(sorted, result{key: k, value: v})
	}
	c.mu.Unlock()
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].value > sorted[j].value
	})
	fmt.Printf("results for [%s]:\n", fileName)
	l := 10
	if l > len(sorted) {
		l = len(sorted)
	}
	for i := 0; i < l; i++ {
		fmt.Printf("  %s:%d\n", sorted[i].key, sorted[i].value)
	}
}

func (c *Cache) ShutDownGracefully(wg *sync.WaitGroup) {
	defer wg.Done()
	c.doneChan <- struct{}{}
	for _, f := range c.m {
		f.wg.Wait()
	}
	fmt.Printf("cache output stopped\n")
}

func (c *Cache) TempPrintMaps() {
	for k1, v1 := range c.m {
		fmt.Printf("MAP[%s] ", k1)
		for k2, v2 := range v1.bow {
			fmt.Printf("%s:%d ", k2, v2)
		}
		fmt.Printf("\n")
	}
}

func (c *Cache) TempPrintMapNames() {
	for k1 := range c.m {
		fmt.Println(k1)
	}
}
