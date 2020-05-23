package storage

import (
    "time"
    "math/rand"

    "github.com/jls83/crisgo/types"
)


type ResultStorage interface {
    Close() (err error)
    GetResultMapKey() types.ResKey
    GetValue(k types.ResKey) (types.ResValue, bool)
    InsertValue(v types.ResValue) types.ResKey
    // TODO: Add explicit `SetValue` method & endpoint for testing
}

// Section: Helper Methods
// This was cribbed from https://www.calhoun.io/creating-random-strings-in-go/
const charset = "abcdefghijklmnopqrstuvwxyz" +
                "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
                "0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
    byteArray := make([]byte, length)
    for i:= range byteArray {
        byteArray[i] = charset[seededRand.Intn(len(charset))]
    }

    return string(byteArray)
}
