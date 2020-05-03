package main

import (
    "fmt"
    "os"
    "strconv"
    "net/http"
    "math/rand"
    "time"

    "crisgo/handlers"
)

func getPortNumberStartMessage(rawPortNumberStr string) (string, string) {
    PORT_NUMBER_MIN := uint64(1)
    PORT_NUMBER_MAX := uint64(65535)

    portNumberAsStr := "8080"

    portNumberAsInt, err := strconv.ParseUint(rawPortNumberStr, 0, 16)

    if err != nil {
        return portNumberAsStr, fmt.Sprintf("There was an error converting %s; starting on %s", rawPortNumberStr, portNumberAsStr)
    } else if (portNumberAsInt < PORT_NUMBER_MIN) || (portNumberAsInt >= PORT_NUMBER_MAX) {
        return portNumberAsStr, fmt.Sprintf("Port %s is out of range; starting on %s", rawPortNumberStr, portNumberAsStr)
    }
    return rawPortNumberStr, fmt.Sprintf("Listening on port %s", portNumberAsStr)
}

type ResultStorage interface {
    Close() (err error)
    GetResultMapKey() handlers.ResKey
    GetValue(k handlers.ResKey) (handlers.ResValue, bool)
    InsertValue(v handlers.ResValue) handlers.ResKey
}

type localStorage struct {
    _innerStorage handlers.ResMap
}

func (s localStorage) Close() (err error) {
    // Since this is just in-memory, don't actually do anything
    return
}

func (s localStorage) GetResultMapKey() handlers.ResKey {
    // FOR NOW
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    return handlers.ResKey(r1.Intn(100))
}

func (s localStorage) GetValue(k handlers.ResKey) (handlers.ResValue, bool) {
    // Get value in _innerStorage
    value, found := s._innerStorage[k]
    return value, found
}

func (s localStorage) InsertValue(v handlers.ResValue) handlers.ResKey {
    // Insert the value into _innerStorage, return the key
    // TODO: Add some error handling; I bet shit can get weird
    var resultKey handlers.ResKey
    hasKey := true

    // Loop until we have a good key
    for hasKey {
        resultKey = s.GetResultMapKey()
        _, hasKey = s._innerStorage[resultKey]
        fmt.Println(hasKey)
    }

    s._innerStorage[resultKey] = v

    return resultKey
}

func main() {
    // TODO: I'm sure there's more shit I can do to lock this down, but...no.
    portNumber := "8080"
    startMessage := fmt.Sprintf("Listening on port %s", portNumber)

    // If the user has passed in a port number arg, check it
    if len(os.Args) >= 2 {
        portNumber, startMessage = getPortNumberStartMessage(os.Args[1])
    }

    m := handlers.ResMap{}

    // Using the `buildFoo` methods allows us to dynamically inject the `resultMap`
    http.HandleFunc("/lengthen/", handlers.BuildLengthen(m))
    http.HandleFunc("/shorten/", handlers.BuildShorten(m))
    http.HandleFunc("/redirector/", handlers.BuildRedirector(m))

    fmt.Println(startMessage)
    addr := ":" + portNumber
    http.ListenAndServe(addr, nil)
}

