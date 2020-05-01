package main

import (
    "fmt"
    "strings"
    "encoding/json"
    "math/rand"
    "net/http"
)

type resKey string
type resValue string
type resultMap map[resKey]resValue

func getResultMapKey() resKey {
    // FOR NOW
    return resKey(rand.Intn(100))
}

func buildLengthen(m resultMap) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json; charset=utf-8")

        // Split out the requested item, then parse & cast it to a `resKey`
        requestedItem := resKey(strings.SplitN(r.URL.Path, "/", 3)[2])

        // Read the item at the hashed address
        // TODO: Use boolean "found" value to return the appropriate HTTP code
        resultValue, _ := m[requestedItem]

        json.NewEncoder(w).Encode(map[string]interface{}{
            "requestedItem": requestedItem,
            "value": resultValue,
        })
    }
}

func buildShorten(m resultMap) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        // Parse form; print to console if we blow up
        err := r.ParseForm()
        if err != nil {
            fmt.Println("ParseForm() err: %v", err)
            return
        }

        // Split out the incomingValue & generate a resultKey
        incomingValue := resValue(r.FormValue("value"))
        resultKey := getResultMapKey()

        // Read the value into the main map
        m[resultKey] = incomingValue

        w.Header().Set("Content-Type", "application/json; charset=utf-8")

        // TODO: Return a response with a CREATED code
        json.NewEncoder(w).Encode(map[string]interface{}{
            "value": incomingValue,
            "location": resultKey,
        })
    }
}

func main() {
    PORT := ":8080"

    fmt.Println("Starting")

    m := resultMap{}

    // Using the `buildFoo` methods allows us to dynamically inject the `resultMap`
    http.HandleFunc("/lengthen/", buildLengthen(m))
    http.HandleFunc("/shorten/", buildShorten(m))

    http.ListenAndServe(PORT, nil)
}

