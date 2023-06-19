package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var (
	configPath string
)

type Config struct {
	ListenAddr string `json:"listen"`
	RemoteAddr string `json:"remote"`
	Password   string `json:"password"`
}

func init() {
	home, _ := homedir.Dir()
	configFilename := ".Websocks.json"
	if len(os.Args) == 2 {
		configFilename = os.Args[1]
	}
	configPath = path.Join(home, configFilename)
}

func (config *Config) SaveConfig() {
	configJson, _ := json.MarshalIndent(config, "", "	")
	err := ioutil.WriteFile(configPath, configJson, 0644)
	if err != nil {
		fmt.Errorf("Fail to save to config file %s: %s", configPath, err)
	}
	log.Printf("Save to config file %s \n", configPath)
}

func (config *Config) ReadConfig() {
	// If config file exist, then start reading data to config instance
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		log.Printf("Read config from %s \n", configPath)
		file, err := os.Open(configPath)
		if err != nil {
			log.Fatalf("Fail to open %s: %s", configPath, err)
		}
		defer file.Close()

		err = json.NewDecoder(file).Decode(config)
		if err != nil {
			log.Fatalf("Invalid JSON config:\n%s", file.Name())
		}
	}
}
