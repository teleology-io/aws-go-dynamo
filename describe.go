package dynamo

import (
	"fmt"
	"reflect"
	"strings"
)

type tableGsi struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type tableSchema struct {
	Table      string     `json:"name"`
	PrimaryKey tableGsi   `json:"pk"`
	Indexes    []tableGsi `json:"gsi"`
}

func schemaFromReflection(tableName string, v interface{}) tableSchema {
	schema := tableSchema{
		Table: tableName,
		PrimaryKey: tableGsi{
			Key:  "",
			Name: "",
		},
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
			schema.PrimaryKey.Key = field.Tag.Get("json")
			schema.PrimaryKey.Name = field.Name
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
