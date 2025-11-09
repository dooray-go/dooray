# Dooray! API in Go

This is the Dooray Client Library for Go.

This library is based on the [Dooray! API Documentation](https://helpdesk.dooray.com/share/pages/9wWo-xwiR66BO5LGshgVTg/2939987647631384419).

## Installing
### go get
```
$ go get -u github.com/dooray-go/dooray
```

## Features

| Category | Feature | Method | Description |
|----------|---------|--------|-------------|
| **Messenger** | Webhook | `PostWebhook` | Send messages via webhook |
| | Direct Send | `DirectSend` | Send direct messages to users |
| **Project** | Get Projects | `GetProjects` | Retrieve list of projects |
| | Get Posts | `GetPosts` | Retrieve posts from a project |
| | Create Post | `CreatePost` | Create a new post in a project |
| **Calendar** | Get Calendars | `GetCalendars` | Retrieve list of calendars |
| | Get Events | `GetEvents` | Retrieve events from calendars |
| | Create Event | `CreateEvent` | Create a new calendar event |

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

## OpenApi Examples

### Get Posts from a Project
```go
package main

import (
    "fmt"
    "log"

    "github.com/dooray-go/dooray/openapi/project"
)

func main() {
    // Create a project client
    projectClient := project.NewDefaultProject()

    // Get posts from a project
    projectID := "your-project-id"
    toMemberIds := "member-id-1,member-id-2" // Optional: filter by member IDs
    postWorkflowClasses := "registered,working" // Optional: filter by workflow classes

    response, err := projectClient.GetPosts("your-dooray-api-key", projectID, toMemberIds, postWorkflowClasses)
    if err != nil {
        log.Fatalf("Failed to get posts: %s", err)
    }

    fmt.Printf("Total posts: %d\n", response.TotalCount)
    for _, post := range response.Result {
        fmt.Printf("Post #%d: %s\n", post.Number, post.Subject)
    }
}
```

### Create a Post in a Project
```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/dooray-go/dooray/openapi/project"
    model "github.com/dooray-go/dooray/openapi/model/project"
    "github.com/dooray-go/dooray/utils"
)

func main() {
    // Create a project client
    projectClient := project.NewDefaultProject()

    // Set due date (24 hours from now)
    dueDate := utils.JsonTime(time.Now().Add(24 * time.Hour))

    // Create a post request
    postRequest := model.PostRequest{
        Subject: "New Task",
        Body: model.PostBody{
            MimeType: "text/html",
            Content:  "<p>This is a new task created via API</p>",
        },
        Users: &model.PostUsers{
            To: []model.PostRecipient{
                {
                    Type: "member",
                    Member: &model.PostMember{
                        OrganizationMemberID: "member-id",
                    },
                },
            },
        },
        DueDate:  &dueDate,
        Priority: "normal", // urgent | high | normal | low
        TagIDs:   []string{"tag-id-1", "tag-id-2"},
    }

    // Create the post
    projectID := "your-project-id"
    response, err := projectClient.CreatePost("your-dooray-api-key", projectID, postRequest)
    if err != nil {
        log.Fatalf("Failed to create post: %s", err)
    }

    fmt.Printf("Post created successfully! ID: %s\n", response.Result.ID)
}
```