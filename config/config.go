package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"path/filepath"
	"sync"
)

var (
	cfg  *tomlConfig
	once sync.Once
)

func Config() *tomlConfig {
	once.Do(func() {
		filePath, err := filepath.Abs("./config.toml")
		if err != nil {
			panic(err)
		}
		fmt.Printf("parse toml file once. filePath: %s\n", filePath)
		if _, err := toml.DecodeFile(filePath, &cfg); err != nil {
			panic(err)
		}
	})
	return cfg
}
