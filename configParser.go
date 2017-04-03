package main

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/kr/pretty"
)

// ParseConfig opens the file at configFileName and unmarshals it into an AppConfig.
func ParseConfig(configFileName string) (*AppConfig, error) {
	file, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Errorf("Error reading configuration file [%v]: [%v]", configFileName, err.Error())
		return nil, err
	}

	var conf AppConfig
	err = json.Unmarshal(file, &conf)
	if err != nil {
		log.Errorf("Error unmarshalling configuration file [%v]: [%v]", configFileName, err.Error())
		return nil, err
	}

	log.Infof("Using configuration: %# v", pretty.Formatter(conf))
	return &conf, nil
}
