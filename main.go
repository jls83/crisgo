package main

import (
    // "errors"
    // "flag"
    "fmt"
    "strconv"
    "strings"
    "encoding/json"
    "net/http"

    "github.com/jls83/crisgo/config"
    "github.com/jls83/crisgo/storage"
    "github.com/jls83/crisgo/types"
)

// Section: Handlers
type LengthenResult struct {
    RequestedItem   types.ResKey      `json:"requestedItem"`
    Value           types.ResValue    `json:"value"`
}

type ShortenResult struct {
    Location    types.ResKey      `json:"location"`
    Value       types.ResValue    `json:"value"`
}

func getRequestedItem(r *http.Request) types.ResKey {
    return types.ResKey(strings.SplitN(r.URL.Path, "/", 3)[2])
}

func BuildRedirector(m storage.ResultStorage) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        requestedItem := getRequestedItem(r)
        resultValue, hasValue := m.GetValue(requestedItem)

        if hasValue {
            http.Redirect(w, r, string(resultValue), http.StatusMovedPermanently)
        } else {
            http.NotFound(w, r)
        }
    }
}

func BuildLengthen(m storage.ResultStorage) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        requestedItem := getRequestedItem(r)

        resultValue, _ := m.GetValue(requestedItem)

        w.Header().Set("Content-Type", "application/json; charset=utf-8")

        result := LengthenResult{
            RequestedItem: requestedItem,
            Value: resultValue,
        }
        json.NewEncoder(w).Encode(result)
    }
}

func BuildShorten(m storage.ResultStorage) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        // Parse form; print to console if we blow up
        err := r.ParseForm()
        if err != nil {
            fmt.Println("ParseForm() err: %v", err)
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("There was an error parsing your link shortening request."))
            return
        }

        incomingValue := types.ResValue(r.FormValue("value"))
        resultKey := m.InsertValue(incomingValue)

        result := ShortenResult{
            Location: resultKey,
            Value: incomingValue,
        }

        w.WriteHeader(http.StatusCreated)
        w.Header().Set("Content-Type", "application/json; charset=utf-8")

        json.NewEncoder(w).Encode(result)
    }
}

// Main
func main() {
    configFile := "base_config.yaml"

    appConfig := config.NewCrisgoConfig(configFile)

    // TODO: Re-add our `flag` items

    // TODO: Have the `New` methods simply take in a `CrisgoConfig` instance!
    // m := storage.NewLocalStorage()
    m := storage.NewSqliteStorage(appConfig)
    defer m.Close()

    err := m.CreateTable()
    if err != nil {
        panic(err)
    }

    // Using the `buildFoo` methods allows us to dynamically inject the ResultStorage instance
    http.HandleFunc("/lengthen/", BuildLengthen(m))
    http.HandleFunc("/shorten/", BuildShorten(m))
    http.HandleFunc("/redirector/", BuildRedirector(m))

    // TODO: Put this out via a logger instead
    fmt.Println("Listening on port", appConfig.PortNumber)
    addr := ":" + strconv.Itoa(appConfig.PortNumber)
    http.ListenAndServe(addr, nil)
}

