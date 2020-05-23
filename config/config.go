package config

import (
    "log"
    "io/ioutil"

    "gopkg.in/yaml.v2"
)

type CrisgoConfig struct {
    Tablename           string
    DatabaseFilePath    string
    PortNumber          int
}


func NewCrisgoConfig(filepath string) *CrisgoConfig {
    fileContents, err := ioutil.ReadFile(filepath)
    if err != nil {
        log.Fatal(err)
    }

    config := CrisgoConfig{}

    err = yaml.Unmarshal(fileContents, &config)
    if err != nil {
        log.Fatal(err)
    }

    return &config
}
