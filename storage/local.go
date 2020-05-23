package storage

import (
    "fmt"

    "github.com/jls83/crisgo/types"
)

// Section: LocalStorage
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
