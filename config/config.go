package config

import (
    "fmt"
    "log"
    "io/ioutil"

    "gopkg.in/yaml.v2"
)

type CrisgoConfig struct {
    Tablename           string  `yaml:"tablename"`
    DatabaseFilePath    string  `yaml:"database_file_path"`
    PortNumber          int     `yaml:"port_number"`
}


func NewCrisgoConfig(filepath string) *CrisgoConfig {
    config := CrisgoConfig{}

    fileContents, err := ioutil.ReadFile(filepath)
    if err != nil {
        // If we can't open the file for whatever reason, simply return the empty struct
        // TODO: This might suck a bit
        fmt.Printf("Error opening", filepath)
        return &config
    }


    err = yaml.Unmarshal(fileContents, &config)
    if err != nil {
        log.Fatal(err)
    }

    return &config
}
