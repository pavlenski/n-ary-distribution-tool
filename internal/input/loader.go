package input

import (
	"fmt"
	"github.com/pavlenski/n-ary-distribution-tool/internal/cruncher"
	"io/ioutil"
	"strings"
	"sync"
)

type job struct {
	filePath  string
	crunchers []chan<- *cruncher.Data
}

type discPool struct {
	jobChan  chan *job
	doneChan chan struct{}
}

type FileLoader struct {
	pool map[string]*discPool
	wg   sync.WaitGroup
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
		wg:   sync.WaitGroup{},
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
			data, err := ioutil.ReadFile(j.filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fileNameSplit := strings.Split(j.filePath, "/")
			fileName := fileNameSplit[len(fileNameSplit)-1]
			for _, c := range j.crunchers {
				c <- cruncher.NewCruncherData(fileName, &data)
			}
			fmt.Printf("loaded file [%s] from disc [%s] size [%dMB]\n", j.filePath, disc, len(data)/1000000)
		case <-doneChan:
			l.wg.Done()
			break
			// should add a sleep channel case where no file input component is working on the specific disc..
			// which means the goroutine for loading files on that disc should be asleep.
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
	fmt.Printf("disc loader was shut down\n")
}
