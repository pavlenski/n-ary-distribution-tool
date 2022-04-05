package app

import (
	"encoding/json"
	"fmt"
	"github.com/pavlenski/n-ary-distribution-tool/internal/input"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const configPath = "./config.json"

type config struct {
	FileInputSleepTime uint64 `json:"file_input_sleep_time,omitempty"`
	Discs              string `json:"discs,omitempty"`
	CounterDataLimit   uint64 `json:"counter_data_limit,omitempty"`
	SortProgressLimit  uint64 `json:"sort_progress_limit,omitempty"`
}

type App struct {
	discs              map[int]string
	fileInputSleepTime time.Duration
	counterDataLimit   uint64
	sortProgressLimit  uint64

	fileLoader      *input.FileLoader
	inputComponents map[string]*input.FileInput
}

func NewApp() *App {
	return &App{
		inputComponents: make(map[string]*input.FileInput),
		discs:           make(map[int]string),
	}
}

func (a *App) Start() {
	fmt.Println("- - - - - n-ary-distribution-tool - - - - -")
	fmt.Println("- - - - - - p.galantic rn3817 - - - - - - -")
	fmt.Printf("- - - - - - - - - %v - - - - - - - - -\n", time.Now().Format("3:04-PM"))
	a.loadConfig()
	a.fileLoader.Run()
	a.run()
}

// loadConfig must be called before running the app & the fileLoader's Run() method
func (a *App) loadConfig() {
	f, err := os.Open(configPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}
	cfg := &config{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		log.Fatalln(err)
	}
	a.configure(cfg)
}

// configure sets up all app variables
func (a *App) configure(cfg *config) {
	ds := strings.Split(cfg.Discs, ";")
	for i, d := range ds {
		a.discs[i+1] = d
	}
	st, err := time.ParseDuration(fmt.Sprintf("%dms", cfg.FileInputSleepTime))
	if err != nil {
		log.Fatalln(err)
	}
	a.fileInputSleepTime = st
	a.counterDataLimit = cfg.CounterDataLimit
	a.sortProgressLimit = cfg.SortProgressLimit
	a.fileLoader = input.NewFileLoader(ds)
}
