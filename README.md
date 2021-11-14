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
	ID     string `json:"id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Expiry int64  `json:"expiry"`
}

// Required!
func (e Expiration) New() interface{} {
	return Expiration{}
}

// Required!
func (e Expiration) PrimaryKey(v interface{}) string {
	var expires Expiration
	mapstructure.Decode(v, expires)
	return expires.ID
}

func main() {
	config := dynamo.Config{
		Table:   "sample-table",
		Log:     false,
		Creater: Expiration{},
	}

	ddb := dynamo.New(config, nil)
}

// Optional 
func Decode(in interface{}, out interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  &out,
	})
	if err != nil {
		return err
	}

	err = decoder.Decode(in)
	if err != nil {
		return err
	}

	return nil
}
```

# Usage 

### Create

```golang
// Does collision checking
createdItem, err := ddb.Create(Expiration{
  ID:     "test-id",
  Email:  "someemail@domain.com",
  Role:   "N/A",
  Expiry: time.Now().Unix() + 300,
})
```

### Update 
```golang
// Does a get, merges updates and writes
updatedItem, err := ddb.Update(Expiration{
  ID:     "test-id",
  Expiry: 0,
})
```

### Put 
```golang
// Does not collision check and just writes
putItem, err := ddb.Put(Expiration{
  ID:     "test-id-2",
  Email:  "anotherone@domain.com",
  Role:   "N/A",
  Expiry: time.Now().Unix(),
})
```

### Get 
```golang
// Uses primary key to get item
gotItem, err := ddb.Get("test-id-2")
```

### Delete 
```golang
// Uses primary key to delete item
err := ddb.Delete("test-id-2")
```

### Query 
```golang
// Can search across n+1 global secondary indexes - must map to data def
results, err := ddb.Query([]dynamo.QueryParams{
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