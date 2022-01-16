# aws-go-dynamo
A wrapper around the aws-sdk dynamo db sdk

# Installation
```
go get github.com/teleology-io/aws-go-dynamo
```

# Configuration
Not all configuration is required, the following example uses mapstructure to convert dynamo responses to our expected structs. Generally you'll want to create a struct representation of your table complete with json mapping, and then implement some ObjectResolver interface methods. 

```golang
package main

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	dynamo "github.com/teleology-io/aws-go-dynamo"
)

type Expiration struct {
	ID     string `json:"id" dynamo:"pk"`
	Email  string `json:"email" dynamo:"gsi,EmailGSI"`
	Role   string `json:"role" dynamo:"gsi,RoleGSI"`
	Expiry int64  `json:"expiry"`
}

func main() {
	sess := session.Must(session.NewSession())
	dd := dynamo.New(sess, &aws.Config{
		Region: aws.String("us-east-1"),
	})

	table := dd.Table(dynamo.Config{
		Table: "sample-table",
		Model: Expiration{},
	})

	
	exp, err := table.Get("123")

	updates := Expiration{
		Role: "admin",
	}

	updatedExp, err := table.Update("123", updates)
}
```

# Usage 

### Create

```golang
// Does collision checking
createdItem, err := table.Create("test-id", Expiration{
  ID:     "test-id",
  Email:  "someemail@domain.com",
  Role:   "N/A",
  Expiry: time.Now().Unix() + 300,
})
```

### Update 
```golang
// Does a get, merges updates and writes
updatedItem, err := table.Update("test-id", Expiration{
  ID:     "test-id",
  Expiry: 0,
})
```

### Put 
```golang
// Does not collision check and just writes
putItem, err := table.Put(Expiration{
  ID:     "test-id-2",
  Email:  "anotherone@domain.com",
  Role:   "N/A",
  Expiry: time.Now().Unix(),
})
```

### Get 
```golang
// Uses primary key to get item
gotItem, err := table.Get("test-id-2")
```

### Delete 
```golang
// Uses primary key to delete item
err := table.Delete("test-id-2")
```

### Query 
```golang
// Can search across n+1 global secondary indexes - must map to data def
results, err := table.Query([]dynamo.QueryParams{
  {
    Key:   "email",
    Value: "anotherone@domain.com",
  },
  {
    Key:   "role",
    Value: "admin",
  },
})
```

# Changelog

**1.0.0**
- Initial Port from @teleology/dynamo