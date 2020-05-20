package storage

import (
    "fmt"
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

type LocalStorage struct {
    _innerStorage types.ResMap
}

func NewLocalStorage() *LocalStorage {
    localStorage := types.ResMap{}
    return &LocalStorage{localStorage}
}

func (s LocalStorage) Close() (err error) {
    // Since this is just in-memory, don't actually do anything
    return
}

// TODO: God help me for using global variables
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

func (s LocalStorage) GetResultMapKey() types.ResKey {
    return types.ResKey(StringWithCharset(5, charset))
}

func (s LocalStorage) GetValue(k types.ResKey) (types.ResValue, bool) {
    value, found := s._innerStorage[k]
    return value, found
}

func (s *LocalStorage) InsertValue(v types.ResValue) types.ResKey {
    // TODO: Add some error handling; I bet shit can get weird
    var resultKey types.ResKey
    hasKey := true

    // Loop until we have a key not already in the map
    for hasKey {
        resultKey = s.GetResultMapKey()
        _, hasKey = s._innerStorage[resultKey]
        fmt.Println(hasKey)
    }

    s._innerStorage[resultKey] = v

    return resultKey
}
