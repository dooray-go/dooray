# Dooray! API in Go

This is the Dooray Client Library for Go.

## Installing
### go get 
```
$ go get -u github.com/dooray-go/dooray
```

## Messenger WebHook Example
```go
package main

import (
    "context"
    "github.com/dooray-go/dooray"
    "log"
    "time"
)

func main() {
    ctx1 := context.Background()
    subCtx1, _ := context.WithTimeout(ctx1, 3*time.Second)
    doorayErr := dooray.PostWebhookContext(subCtx1, "[Your WebHook URL]", &dooray.WebhookMessage{
        BotName: "dooray-go",
        Text:    "Hello",
    })
    
    if doorayErr != nil {
        log.Printf("dial error: %s", doorayErr.Error())
    }
}
```

## OpenApi Example
```go
package main

import (
    "context"
    "github.com/dooray-go/dooray"
    "log"
    "time"
)

func main() {
	projectType := "public"
	scope := "private"
	state := "active"

	response, err := NewDefaultProject().GetProjects("{dooray-api-key}", projectType, scope, state)
	if err != nil {
		t.Errorf("Expected not to receive error: %s", err)
	}

	fmt.Println(response)
}
```