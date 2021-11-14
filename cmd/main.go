package main

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	dynamo "github.com/teleology-io/aws-go-dynamo"
)

type Expiration struct {
	ID     string `json:"id"`
	Expiry int64  `json:"expiry"`
}

func (e Expiration) New() interface{} {
	return Expiration{}
}

func (e Expiration) PrimaryKey(v interface{}) string {
	var expires Expiration
	mapstructure.Decode(v, expires)
	return expires.ID
}

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

func main() {
	config := dynamo.Config{
		Table:   "sample-table",
		Log:     false,
		Creater: Expiration{},
	}

	ddb := dynamo.New(config, nil)

	created, err := ddb.Create(Expiration{
		ID:     "test-id",
		Expiry: time.Now().Unix() + 300,
	})
	if err != nil {
		// handle err
	}

	var createdExpiration Expiration
	err = Decode(created, &createdExpiration)
	if err != nil {
		// handle err
	}

	fmt.Println(createdExpiration.Expiry)

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
}
