package main

import (
    "fmt"
    "os"
    "strconv"
    "strings"
    "time"
    "encoding/json"
    "math/rand"
    "net/http"
)

// Section: Types
type ResKey string
type ResValue string
type ResMap map[ResKey]ResValue

// Section: Storage
type ResultStorage interface {
    Close() (err error)
    GetResultMapKey() ResKey
    GetValue(k ResKey) (ResValue, bool)
    InsertValue(v ResValue) ResKey
}

type LocalStorage struct {
    _innerStorage ResMap
}

func (s LocalStorage) Close() (err error) {
    // Since this is just in-memory, don't actually do anything
    return
}

func NewLocalStorage() *LocalStorage {
    localStorage := ResMap{}
    return &LocalStorage{localStorage}
}

func (s *LocalStorage) GetResultMapKey() ResKey {
    // FOR NOW
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    return ResKey(r1.Intn(100))
}

func (s *LocalStorage) GetValue(k ResKey) (ResValue, bool) {
    // Get value in _innerStorage
    value, found := s._innerStorage[k]
    return value, found
}

func (s *LocalStorage) InsertValue(v ResValue) ResKey {
    // Insert the value into _innerStorage, return the key
    // TODO: Add some error handling; I bet shit can get weird
    var resultKey ResKey
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

// Section: Handlers
func getResultMapKey() ResKey {
    // FOR NOW
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    return ResKey(r1.Intn(100))
}

func getRequestedItem(r *http.Request) ResKey {
    return ResKey(strings.SplitN(r.URL.Path, "/", 3)[2])
}

func BuildRedirector(m *LocalStorage) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        requestedItem := getRequestedItem(r)
        resultValue, hasValue := m.GetValue(requestedItem)

        if hasValue {
            http.Redirect(w, r, string(resultValue), http.StatusMovedPermanently)
        } else {
            http.NotFound(w, r)
        }
    }
}

func BuildLengthen(m *LocalStorage) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        requestedItem := getRequestedItem(r)

        // Read the item at the hashed address
        // TODO: Use boolean "found" value to return the appropriate HTTP code
        resultValue, _ := m.GetValue(requestedItem)

        w.Header().Set("Content-Type", "application/json; charset=utf-8")

        json.NewEncoder(w).Encode(map[string]interface{}{
            "requestedItem": requestedItem,
            "value": resultValue,
        })
    }
}

func BuildShorten(m *LocalStorage) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        // Parse form; print to console if we blow up
        err := r.ParseForm()
        if err != nil {
            fmt.Println("ParseForm() err: %v", err)
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("There was an error parsing your link shortening request."))
            return
        }

        incomingValue := ResValue(r.FormValue("value"))
        resultKey := m.InsertValue(incomingValue)

        w.WriteHeader(http.StatusCreated)
        w.Header().Set("Content-Type", "application/json; charset=utf-8")

        json.NewEncoder(w).Encode(map[string]interface{}{
            "value": incomingValue,
            "location": resultKey,
        })
    }
}

// Section: Other
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

// Main
func main() {
    // TODO: I'm sure there's more shit I can do to lock this down, but...no.
    portNumber := "8080"
    startMessage := fmt.Sprintf("Listening on port %s", portNumber)

    // If the user has passed in a port number arg, check it
    if len(os.Args) >= 2 {
        portNumber, startMessage = getPortNumberStartMessage(os.Args[1])
    }

    m := NewLocalStorage()
    defer m.Close()

    // Using the `buildFoo` methods allows us to dynamically inject the `resultMap`
    http.HandleFunc("/lengthen/", BuildLengthen(m))
    http.HandleFunc("/shorten/", BuildShorten(m))
    http.HandleFunc("/redirector/", BuildRedirector(m))

    fmt.Println(startMessage)
    addr := ":" + portNumber
    http.ListenAndServe(addr, nil)
}

