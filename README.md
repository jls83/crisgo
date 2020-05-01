# Crisgo
_It's shortening!_

## Quickstart
```shell
go build && ./crisgo
```

Then, in another shell:
```shell
# This should return a response with `"hey"`
curl http://localhost:8080/lengthen/1

# This should return a response with the address for the shortened item
curl -d "value=something" http://localhost:8080/shorten/
```
