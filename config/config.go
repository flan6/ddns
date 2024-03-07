package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type Config struct {
	UpdateInterval int `json:"update_interval"`
}

var (
	config *Config
	mutex  sync.RWMutex
)

func loadConfig(path string) error {
	// Read and unmarshal config
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	tempConfig := Config{}
	if err := json.Unmarshal(file, &tempConfig); err != nil {
		return err
	}

	config = &tempConfig
	return nil
}

func GetConfig() *Config {
	mutex.RLock()
	defer mutex.RUnlock()
	return config
}

func WatchConfig(path string) {
	// Initial load
	if err := loadConfig(path); err != nil {
		panic(err)
	}

	lastModTime := time.Now()
	for {
		// Check for modification
		stat, err := os.Stat(path)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		if stat.ModTime().After(lastModTime) {
			if err := loadConfig(path); err != nil {
				fmt.Println("could not update config", err)
			} else {
				lastModTime = stat.ModTime()
			}
		}

		time.Sleep(60 * time.Second)
	}
}
