package input

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

func (i *FileInput) checkDisk() {
	fmt.Printf("input [%s] started working\n", i.Name)
	i.working = true
	go i.crawl()
}

func (i *FileInput) crawl() {
	for _, dir := range i.directories {
		err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() || !strings.HasSuffix(info.Name(), ".txt") {
				return nil
			}
			if !i.recentlyModified[info.Name()].Before(info.ModTime()) {
				fmt.Printf("file [%s] was recently modified, skipping\n", path)
				return nil
			}
			i.recentlyModified[info.Name()] = time.Now()
			time.Sleep(2 * time.Second)
			fmt.Printf("path [%s] info [%s]\n", path, info.Name())
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Printf("input [%s] finished working\n", i.Name)
	// in case the file input component got paused during its work
	// we do not wish to snooze, but just stay paused
	if i.state != Paused {
		i.setState(Paused)
		go i.snooze()
	} else {
		// if the component is paused, we will instantly cancel the wait for the snooze func
		i.wg.Done()
	}

	i.working = false
	i.wg.Done()
}

func (i *FileInput) loadAndSendFile(filePath string) {
	i.poolChan <- struct{}{}
	// load the file & send
	<-i.poolChan
}
