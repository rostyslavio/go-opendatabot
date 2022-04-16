# Go Opendatabot
Go mod for [Opendatabot](https://opendatabot.ua/en) service.

## Installation
```shell
go get -u github.com/rostyslavio/go-opendatabot
```

## Documentation
https://docs.opendatabot.com/

## Examples
```go
package main

import (
	odb "github.com/rostyslavio/go-opendatabot"
	"log"
)

func main() {
	s, _ := odb.NewOdbClient(odb.WithApiKey("API_KEY"))

	result, err := s.GetGovernmentCompany("14360570")

	if err != nil {
		log.Fatal(err);
	}

	log.Println(result);
}
```


## Licence
This package is licensed under the MIT license. See LICENSE for details.
