package dynamo

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type QueryParams struct {
	Key   string
	Value interface{}
}

type queryExpression struct {
	IndexName  string
	Expression string
}

var PUT_RETURN_VALUES = "ALL_OLD"

func buildPrimaryKey(t TableDef, pk string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		t.baseParams.Key: {
			S: aws.String(pk),
		},
	}
}

func merge(maps ...map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for _, it := range maps {
		for k, v := range it {
			if !reflect.ValueOf(v).IsZero() {
				res[k] = v
			}
		}
	}
	return res
}

func (t TableDef) Get(pk string) (interface{}, error) {
	params := &dynamodb.GetItemInput{
		TableName: &t.baseParams.Table,
		Key:       buildPrimaryKey(t, pk),
	}
	res, err := db.GetItem(params)
	if logger != nil {
		logger.Println("GET_ITEM: ", params)
	}

	if err != nil {
		return nil, err
	}

	if res.Item == nil {
		return nil, errors.New("No item found with pkey " + pk)
	}

	out := reflect.New(t.t).Interface()
	err = dynamodbattribute.UnmarshalMap(res.Item, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (t TableDef) Put(in interface{}) (interface{}, error) {
	item, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		return nil, err
	}

	params := &dynamodb.PutItemInput{
		TableName:    &t.baseParams.Table,
		Item:         item,
		ReturnValues: &PUT_RETURN_VALUES,
	}
	_, err = db.PutItem(params)
	if logger != nil {
		logger.Println("PUT_ITEM: ", params)
	}

	if err != nil {
		return nil, err
	}

	return in, nil
}

func (t TableDef) Delete(pk string) error {
	params := &dynamodb.DeleteItemInput{
		TableName: &t.baseParams.Table,
		Key:       buildPrimaryKey(t, pk),
	}
	_, err := db.DeleteItem(params)
	if logger != nil {
		logger.Println("DELETE_ITEM: ", params)
	}

	if err != nil {
		return err
	}

	return nil
}

func (t TableDef) Create(pk string, v interface{}) (interface{}, error) {
	params := &dynamodb.GetItemInput{
		TableName: &t.baseParams.Table,
		Key:       buildPrimaryKey(t, pk),
	}
	res, err := db.GetItem(params)
	if logger != nil {
		logger.Println("CREATE_COLLISION_CHECK: ", params)
	}

	if err != nil {
		return nil, err
	}

	// hasn't been created yet
	if res.Item != nil {
		return nil, errors.New("A record already exists for '" + pk + "'")
	}

	return t.Put(v)
}

func (t TableDef) Update(pk string, v interface{}) (interface{}, error) {
	params := &dynamodb.GetItemInput{
		TableName: &t.baseParams.Table,
		Key:       buildPrimaryKey(t, pk),
	}
	res, err := db.GetItem(params)
	if logger != nil {
		logger.Println("UPDATE_COLLISION_CHECK: ", params)
	}
	if err != nil {
		return nil, err
	}

	var exist map[string]interface{}
	err = dynamodbattribute.UnmarshalMap(res.Item, &exist)
	if err != nil {
		return nil, err
	}

	out, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var updates map[string]interface{}
	err = json.Unmarshal(out, &updates)
	if err != nil {
		return nil, err
	}

	merged := merge(exist, updates)
	return t.Put(merged)
}

func (t TableDef) Query(qps []QueryParams) ([]interface{}, error) {
	var expressions []queryExpression
	var attributes []map[string]interface{}

	for _, qp := range qps {
		var exist *tableGsi = nil
		for _, ind := range t.baseParams.Indexes {
			if ind.Key == qp.Key {
				exist = &ind

				attributes = append(attributes, map[string]interface{}{
					":" + qp.Key: qp.Value,
				})

				expressions = append(expressions, queryExpression{
					IndexName:  exist.Name,
					Expression: fmt.Sprintf("%s = :%s", qp.Key, qp.Key),
				})
			}
		}
	}

	first := expressions[0]
	rest := expressions[1:]

	mergedExpressions := merge(attributes...)
	expressionAttributeValues, err := dynamodbattribute.MarshalMap(mergedExpressions)
	if err != nil {
		return nil, err
	}

	params := dynamodb.QueryInput{
		TableName:                 &t.baseParams.Table,
		IndexName:                 &first.IndexName,
		KeyConditionExpression:    &first.Expression,
		ExpressionAttributeValues: expressionAttributeValues,
	}

	if len(rest) > 0 {
		var expStr []string
		for _, ex := range rest {
			expStr = append(expStr, ex.Expression)
		}

		filterExpression := strings.Join(expStr, " AND ")
		params.FilterExpression = &filterExpression
	}

	var items []interface{}
	for ok := true; ok; {
		out, err := db.Query(&params)
		if logger != nil {
			logger.Println("QUERY: ", params)
		}

		if err != nil {
			return nil, err
		}

		if out.Items != nil {
			for _, outItem := range out.Items {
				newp := reflect.New(t.t).Interface()
				err = dynamodbattribute.UnmarshalMap(outItem, &newp)
				if err != nil {
					return nil, err
				}

				items = append(items, newp)
			}
		}

		if out.LastEvaluatedKey != nil {
			params.ExclusiveStartKey = out.LastEvaluatedKey
		} else {
			break
		}
	}

	return items, nil
}
