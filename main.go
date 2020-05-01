package main

import (
    "fmt"
    "os"
    "strings"
    "strconv"
    "time"
    "encoding/json"
    "math/rand"
    "net/http"
)

type resKey string
type resValue string
type resultMap map[resKey]resValue

func getResultMapKey() resKey {
    // FOR NOW
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    return resKey(r1.Intn(100))
}

func buildLengthen(m resultMap) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        // Split out the requested item, then parse & cast it to a `resKey`
        requestedItem := resKey(strings.SplitN(r.URL.Path, "/", 3)[2])

        // Read the item at the hashed address
        // TODO: Use boolean "found" value to return the appropriate HTTP code
        resultValue, _ := m[requestedItem]

        w.Header().Set("Content-Type", "application/json; charset=utf-8")

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

        // Try generating a key for our result map. If there's already a result in place,
        // regenerate the key.
        var resultKey resKey
        hasKey := true

        for hasKey {
            resultKey = getResultMapKey()
            _, hasKey = m[resultKey]
            fmt.Println(hasKey)
        }

        incomingValue := resValue(r.FormValue("value"))
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
    PORT_NUMBER_MIN := uint64(1)
    PORT_NUMBER_MAX := uint64(65535)

    // TODO: I'm sure there's more shit I can do to lock this down, but...no.
    portNumber := "8080"
    startMessage := fmt.Sprintf("Listening on port %s", portNumber)

    if len(os.Args) >= 2 {
        rawPortNumberStr := os.Args[1]
        portNumberAsInt, err := strconv.ParseUint(rawPortNumberStr, 0, 16)

        if err != nil {
            startMessage = fmt.Sprintf("There was an error converting %s; starting on %s",
                                       rawPortNumberStr, portNumber)
        } else if (portNumberAsInt < PORT_NUMBER_MIN) || (portNumberAsInt >= PORT_NUMBER_MAX) {
            startMessage = fmt.Sprintf("Port %s is out of range; starting on %s",
                                       rawPortNumberStr, portNumber)
        } else {
            portNumber = rawPortNumberStr
            startMessage = fmt.Sprintf("Listening on port %s", portNumber)
        }
    }

    m := resultMap{}

    // Using the `buildFoo` methods allows us to dynamically inject the `resultMap`
    http.HandleFunc("/lengthen/", buildLengthen(m))
    http.HandleFunc("/shorten/", buildShorten(m))

    fmt.Println(startMessage)
    addr := ":" + portNumber
    http.ListenAndServe(addr, nil)
}

