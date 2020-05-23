## To Do
1. Add new `storage` package that accounts for the following storage types:
    - Redis-backed
2. Add configuration items
    - YAML parser (via `go-yaml`, I'm not insane)
    - CLI options
        - `config-file`
        - Each option in the YAML should be represented, as we can use it to build up the `Config` object
2. Add (async!) event hooks to the listener methods

## To Consider
* How should we account for "creation date"/"last accessed"/etc?
