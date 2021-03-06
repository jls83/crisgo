package config

import (
    "fmt"
    "log"
    "reflect"
    "io/ioutil"

    "gopkg.in/yaml.v2"
)

const DEFAULT_PORT_NUMBER = 8080

const DEFAULT_SQLITE_FILE_PATH = "crisgo.db"
const DEFAULT_SQLITE_TABLE_NAME = "shortened_urls"

// NOTE: Go v1.13 introduces this method in the `reflect` package directly, but my personal
// machine only has v1.12
func isZero(v interface{}) (bool, error) {
    t := reflect.TypeOf(v)
    if !t.Comparable() {
        return false, fmt.Errorf("type is not comparable: %v", t)
    }
    return v == reflect.Zero(t).Interface(), nil
}

type CrisgoConfig struct {
    TableName           string  `yaml:"tablename"`
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
    res, err := isZero(config.TableName)
    if res || err != nil {
        config.TableName = DEFAULT_SQLITE_TABLE_NAME
    }

    // Get file_path
    res, err = isZero(config.DatabaseFilePath)
    if res || err != nil {
        config.DatabaseFilePath = DEFAULT_SQLITE_FILE_PATH
    }

    // Get port_number
    res, err = isZero(config.PortNumber)
    if res || err != nil {
        config.PortNumber = DEFAULT_PORT_NUMBER
    }

    return &config
}
