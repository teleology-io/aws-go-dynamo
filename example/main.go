package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	dynamo "github.com/teleology-io/aws-go-dynamo"
)

// {
// 	"id": "d259a6cb-ef5d-4130-b960-a9f02ecf8472",
// 	"hid": "c3b551316cec3f13600f46ce4608c390cf604df3",
// 	"event": "user.deleted",
// 	"created": 1641102115318,
// 	"url": "https://foundation-api-dev.teleology.io/webhooks/user-deleted",
// 	"name": "foundation user removed (dev)",
// 	"type": "webhook"
// }

// pipeline-dev-subscriptions

type PipelineSubscriptions struct {
	ID      string `json:"id" dynamo:"pk"`
	Hid     string `json:"hid" dynamo:"gsi,HashGSI"`
	Event   string `json:"event"`
	Created int64  `json:"created"`
	URL     string `json:"url"`
	Name    string `json:"name"`
	Type    string `json:"type"`
}

func (c PipelineSubscriptions) New() interface{} {
	return PipelineSubscriptions{}
}

func (c PipelineSubscriptions) PrimaryKey(it interface{}) string {
	switch it.(type) {
	case PipelineSubscriptions:
		fmt.Println("type is cat")
		return PipelineSubscriptions(it.(PipelineSubscriptions)).ID
	default:
		var c = PipelineSubscriptions{}
		// convert.Decode(it, &c)
		return c.ID
	}
}

func main() {
	// p := PipelineSubscriptions{
	// 	ID:  "123",
	// 	Hid: "231",
	// }

	sess := session.Must(session.NewSession())
	dd := dynamo.New(sess, &aws.Config{
		Region: aws.String("us-east-1"),
	})

	table := dd.Table(dynamo.Config{
		Table: "pipeline-dev-subscriptions",
		Model: PipelineSubscriptions{},
	})

	sub, err := table.Get("d259a6cb-ef5d-4130-b960-a9f02ecf8472")

	fmt.Println(sub, err)
}
