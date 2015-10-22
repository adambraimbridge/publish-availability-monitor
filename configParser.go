package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ParseConfig(configFileName string) (*AppConfig, error) {
	file, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return nil, err
	}
	var conf AppConfig
	err = json.Unmarshal(file, &conf)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Config: %v\n", conf)
	return &conf, nil
}
