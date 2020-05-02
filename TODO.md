## To Do
1. Add new `storage` package that accounts for the following storage types:
    - In memory (essentially what we have now)
    - SQL-backed
    - Redis-backed
    ```go
    type resultStorage interface {
        GetResultMapKey() resKey  // For now, a copy/paste job
        GetValue(k resKey) (resValue, bool)  // Returns the value (if any), as well as a bool for if there was a value. Parallel to the `map` access.
        SetValue(k resKey, v resValue) bool  // Returns a bool if the `set` was successful
        InsertValue(v resValue) resKey  // Returns the key where the value can be found
    }
    ```
2. Add new API types
    - "Shorten" input
    - "Lengthen" output

## To Consider
* How should we account for "creation date"/"last accessed"/etc?
