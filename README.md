go-cloudant: A Cloudant API wrapper for Go
===

Supports go 1.2 or later

## Documentation
````sh
(coming soon)
````

## Installation
```sh
go get -u github.com/obieq/go-cloudant
```

## Client
```sh
import "github.com/obieq/go-cloudant"

var client *Client = NewClient("Your Cloudant URI", "Your Cloudant API Key", "Your Cloudant API Password")
```

## Examples
````sh
(coming soon)

Detailed examples can be found in the test files
````
## Tests

100% code coverage via ginkgo and gomega

```sh
go test -cover
      or
ginkgo -cover
```
