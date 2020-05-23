# Crisgo
_It's shortening!_

## Quickstart
```shell
go run main.go
```

Then, in another shell:
```shell
# This should return a response with the address for the shortened item
curl -d "value=something" http://localhost:8080/shorten/

# This should return a response with your shortened URL
curl http://localhost:8080/lengthen/{INSERT_RESULT_HERE}
```

You can also hit the `redirector` endpoint in a browser like so: `http://localhost:8080/redirector/{INSERT_RESULT_HERE}`
