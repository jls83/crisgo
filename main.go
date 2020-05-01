package main

import (
    "fmt"
    "strconv"
    "strings"
    "encoding/json"
    "hash/fnv"
    "net/http"
)

type resKey uint32
type resValue string
type resultMap map[resKey]resValue

func getResultMapKey(s resValue) resKey {
    h := fnv.New32a()
    h.Write([]byte(s))

    return resKey(h.Sum32())
}

func buildLengthen(m resultMap) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json; charset=utf-8")

        // Split out the requested item, then parse & cast it to a `resKey`
        requestedItem := strings.SplitN(r.URL.Path, "/", 3)[2]
        requestedItemAsInt, _ := strconv.Atoi(requestedItem)
        resultKey := resKey(requestedItemAsInt)

        // Read the item at the hashed address
        // TODO: Get element from array
        // TODO: Use boolean "found" value to return the appropriate HTTP code
        resultValue, _ := m[resultKey]

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

        // Hash the incoming `value`
        // TODO: We should salt these as well
        incomingValue := resValue(r.FormValue("value"))
        resultKey := getResultMapKey(incomingValue)

        // Read the value into the main map
        // TODO: Append to an array if we have multiple values
        m[resultKey] = incomingValue

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

    m := resultMap{1: "hey", 2: "ho"}

    // Using the `buildFoo` methods allows us to dynamically inject the `resultMap`
    http.HandleFunc("/lengthen/", buildLengthen(m))
    http.HandleFunc("/shorten/", buildShorten(m))

    http.ListenAndServe(PORT, nil)
}

