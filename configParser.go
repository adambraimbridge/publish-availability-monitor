package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/kr/pretty"
)

func ParseConfig(configFileName string) (*AppConfig, error) {
	file, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Printf("Error reading configuration file [%v]: [%v]", configFileName, err.Error())
		return nil, err
	}

	var conf AppConfig
	err = json.Unmarshal(file, &conf)
	if err != nil {
		log.Printf("Error unmarshalling configuration file [%v]: [$v]", configFileName, err.Error())
		return nil, err
	}

	info.Printf("Using configuration: %# v", pretty.Formatter(conf))
	return &conf, nil
}
