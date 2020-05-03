package handlers

import (
    "fmt"
    "strings"
    "net/http"
    "encoding/json"
    "time"
    "math/rand"
)

type ResKey string
type ResValue string
type ResMap map[ResKey]ResValue

func getResultMapKey() ResKey {
    // FOR NOW
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    return ResKey(r1.Intn(100))
}

func BuildRedirector(m ResMap) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        requestedItem := ResKey(strings.SplitN(r.URL.Path, "/", 3)[2])

        resultValue, hasValue := m[requestedItem]

        if hasValue {
            http.Redirect(w, r, string(resultValue), 301)
            return
        }
        http.NotFound(w, r)
    }
}

func BuildLengthen(m ResMap) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        // Split out the requested item, then parse & cast it to a `ResKey`
        requestedItem := ResKey(strings.SplitN(r.URL.Path, "/", 3)[2])

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

func BuildShorten(m ResMap) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        // Parse form; print to console if we blow up
        err := r.ParseForm()
        if err != nil {
            fmt.Println("ParseForm() err: %v", err)
            return
        }

        // Try generating a key for our result map. If there's already a result in place,
        // regenerate the key.
        var resultKey ResKey
        hasKey := true

        for hasKey {
            resultKey = getResultMapKey()
            _, hasKey = m[resultKey]
            fmt.Println(hasKey)
        }

        incomingValue := ResValue(r.FormValue("value"))
        m[resultKey] = incomingValue

        w.Header().Set("Content-Type", "application/json; charset=utf-8")

        // TODO: Return a response with a CREATED code
        json.NewEncoder(w).Encode(map[string]interface{}{
            "value": incomingValue,
            "location": resultKey,
        })
    }
}
