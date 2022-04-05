package input

import (
	"fmt"
	"github.com/pavlenski/n-ary-distribution-tool/internal/cruncher"
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

// crawl goes through its component's given directories & gathers a group of files which
// will be sent as jobs (each file) to the fileLoader component
func (i *FileInput) crawl() {
	var filePaths []string
	for _, dir := range i.directories {
		err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() || !strings.HasSuffix(info.Name(), ".txt") {
				return nil
			}
			if !i.recentlyModified[path].Before(info.ModTime()) {
				//fmt.Printf("file [%s] wasn't recently modified, skipping\n", path)
				return nil
			}
			i.recentlyModified[path] = time.Now()
			filePaths = append(filePaths, path)
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
	}

	i.createJobs(filePaths)
	i.setState(Paused)
	go i.snooze()
}

// createJob should be a goroutine working since the start of the input lifecycle.
// it should be constantly working and waiting for jobs to work on
func (i *FileInput) createJobs(filePaths []string) {
	for _, filePath := range filePaths {
		if i.state == Paused {
			break
		}
		var crunchers []chan<- *cruncher.Data
		for _, dataChan := range i.crunchers {
			crunchers = append(crunchers, dataChan)
		}
		j := &job{
			filePath:  filePath,
			crunchers: crunchers,
		}
		i.jobChan <- j
	}
}
