package config

import (
    "fmt"
    "log"
    "io/ioutil"

    "gopkg.in/yaml.v2"
)

const DEFAULT_PORT_NUMBER = 8080

const DEFAULT_SQLITE_FILE_PATH = "crisgo.db"
const DEFAULT_SQLITE_TABLE_NAME = "shortened_urls"

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

    // Check our values
    tablename := DEFAULT_SQLITE_TABLE_NAME
    if reflect.ValueOf(config.Tablename).IsZero() {
        config.Tablename = DEFAULT_SQLITE_TABLE_NAME
    }

    // Get file_path
    if reflect.ValueOf(config.DatabaseFilePath).IsZero() {
        config.DatabaseFilePath = DEFAULT_SQLITE_FILE_PATH
    }

    // Get port_number
    if reflect.ValueOf(config.PortNumber).IsZero() {
        config.PortNumber = DEFAULT_PORT_NUMBER
    }

    return &config
}
