package input

import (
	"fmt"
	"sync"
	"time"
)

type job struct {
	filePath string
	// cruncher []*cruncher.Cruncher
}

type discPool struct {
	jobChan  chan *job
	doneChan chan struct{}
}

type FileLoader struct {
	pool map[string]*discPool
	wg   *sync.WaitGroup
}

func newDiscPool() *discPool {
	return &discPool{
		jobChan:  make(chan *job, 1),
		doneChan: make(chan struct{}, 1),
	}
}

func NewFileLoader(discs []string) *FileLoader {
	m := make(map[string]*discPool, len(discs))
	for _, d := range discs {
		m[d] = newDiscPool()
	}
	return &FileLoader{
		pool: m,
		wg:   &sync.WaitGroup{},
	}
}

func (l *FileLoader) Run() {
	for disc := range l.pool {
		l.wg.Add(1)
		go l.LoadAndSendFileFromDisc(disc)
	}
}

func (l *FileLoader) LoadAndSendFileFromDisc(disc string) {
	fmt.Printf("loader listening on disc [%s]\n", disc)
	jobChan := l.pool[disc].jobChan
	doneChan := l.pool[disc].doneChan
	for {
		select {
		case j := <-jobChan:
			time.Sleep(3 * time.Second)
			fmt.Printf("loaded file [%s] from disc [%s]\n", j.filePath, disc)
		case <-doneChan:
			l.wg.Done()
			fmt.Printf("shutting down disc [%s] file loader\n", disc)
			break
		}
	}
}

func (l *FileLoader) GetJobChan(disc string) chan *job {
	if discPool, ok := l.pool[disc]; !ok {
		fmt.Printf("passed disc [%s] doesn't exist, returning nil chan\n", disc)
		return nil
	} else {
		return discPool.jobChan
	}
}

func (l *FileLoader) ShutDownGracefully(inputWg *sync.WaitGroup) {
	defer inputWg.Done()
	for _, discPool := range l.pool {
		discPool.doneChan <- struct{}{}
	}
	l.wg.Wait()
}
