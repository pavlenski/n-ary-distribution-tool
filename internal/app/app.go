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

	inputHandlers map[string]*input.FileInput
}

func NewApp() *App {
	return &App{
		inputHandlers: make(map[string]*input.FileInput),
		discs:         make(map[int]string),
	}
}

func (a *App) Start() {
	a.loadConfig()
	a.run()
}

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
}
