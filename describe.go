package dynamo

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type tableGsi struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type tableSchema struct {
	Table   string     `json:"name"`
	Key     string     `json:"pk"`
	Indexes []tableGsi `json:"gsi"`
}

func schemaFromReflection(tableName string, v interface{}) tableSchema {
	schema := tableSchema{
		Table:   tableName,
		Key:     "",
		Indexes: []tableGsi{},
	}

	t := reflect.TypeOf(v)

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		// Get the field
		field := t.Field(i)

		// Get tag
		tag := field.Tag.Get("dynamo")

		if tag == "" || tag == "-" {
			continue
		}

		if strings.Contains(tag, "pk") {
			schema.Key = field.Tag.Get("json")
		}

		if strings.Contains(tag, "gsi") {
			gsi := strings.Split(tag, ",")
			if len(gsi) == 1 {
				fmt.Println(`gsi tag must include name, ex) gsi,IndexName`)
			}
			schema.Indexes = append(schema.Indexes, tableGsi{
				Key:  field.Tag.Get("json"),
				Name: gsi[1],
			})
		}
	}

	return schema
}

func hash(elements []*dynamodb.KeySchemaElement) string {
	key := ""
	for _, schema := range elements {
		if *schema.KeyType == "HASH" {
			key = *schema.AttributeName
		}
	}

	return key
}

func schemaFromDescribeTable(tableName string) tableSchema {
	// get aws description
	description, err := db.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: &tableName,
	})

	if err != nil {
		panic(err)
	}

	out := tableSchema{
		Table: *description.Table.TableName,
		// Get primary index
		Key: hash(description.Table.KeySchema),
	}

	// Get secondary indexes
	indexes := []tableGsi{}
	for _, secondary := range description.Table.GlobalSecondaryIndexes {
		indexes = append(indexes, tableGsi{
			Key:  hash(secondary.KeySchema),
			Name: *secondary.IndexName,
		})
	}

	out.Indexes = indexes
	return out
}
