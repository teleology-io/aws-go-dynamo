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
	key   string
	value interface{}
}

type queryExpression struct {
	IndexName  string
	Expression string
}

var PUT_RETURN_VALUES = "ALL_OLD"

func key(ds *DynamoService, pk string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		ds.baseParams.key: {
			S: aws.String(pk),
		},
	}
}

func merge(maps ...map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for _, it := range maps {
		for k, v := range it {
			res[k] = v
		}
	}
	return res
}

func (ds *DynamoService) Get(pk string, v interface{}) error {
	params := &dynamodb.GetItemInput{
		TableName: &ds.baseParams.table,
		Key:       key(ds, pk),
	}
	res, err := ds.svc.GetItem(params)
	if ds.logger != nil {
		ds.logger.Println("GET_ITEM: ", params)
	}

	if err != nil {
		return err
	}

	if res.Item == nil {
		return errors.New("No item found with pkey " + pk)
	}

	err = dynamodbattribute.UnmarshalMap(res.Item, &v)
	if err != nil {
		return nil
	}

	return nil
}

func (ds *DynamoService) Put(in interface{}) error {
	item, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		return err
	}

	params := &dynamodb.PutItemInput{
		TableName:    &ds.baseParams.table,
		Item:         item,
		ReturnValues: &PUT_RETURN_VALUES,
	}
	_, err = ds.svc.PutItem(params)
	if ds.logger != nil {
		ds.logger.Println("PUT_ITEM: ", params)
	}

	if err != nil {
		return err
	}

	return nil
}

func (ds *DynamoService) Delete(pk string) (interface{}, error) {
	params := &dynamodb.DeleteItemInput{
		TableName: &ds.baseParams.table,
		Key:       key(ds, pk),
	}
	res, err := ds.svc.DeleteItem(params)
	if ds.logger != nil {
		ds.logger.Println("DELETE_ITEM: ", params)
	}

	if err != nil {
		return nil, err
	}

	var v interface{}
	err = dynamodbattribute.UnmarshalMap(res.Attributes, &v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (ds *DynamoService) Create(pk string, v interface{}) error {
	params := &dynamodb.GetItemInput{
		TableName: &ds.baseParams.table,
		Key:       key(ds, pk),
	}
	res, err := ds.svc.GetItem(params)
	if ds.logger != nil {
		ds.logger.Println("CREATE_COLLISION_CHECK: ", params)
	}

	if err != nil {
		return err
	}

	// hasn't been created yet
	if res.Item != nil {
		return errors.New("A record already exists for '" + pk + "'")
	}

	return ds.Put(v)
}

func (ds *DynamoService) Update(pk string, v interface{}) error {
	params := &dynamodb.GetItemInput{
		TableName: &ds.baseParams.table,
		Key:       key(ds, pk),
	}
	res, err := ds.svc.GetItem(params)
	if ds.logger != nil {
		ds.logger.Println("UPDATE_COLLISION_CHECK: ", params)
	}
	if err != nil {
		return err
	}

	var exist map[string]interface{}
	err = dynamodbattribute.UnmarshalMap(res.Item, &exist)
	if err != nil {
		return nil
	}

	out, err := json.Marshal(v)
	if err != nil {
		return err
	}

	var updates map[string]interface{}
	err = json.Unmarshal(out, &updates)
	if err != nil {
		return err
	}

	for k, v := range exist {
		uk := updates[k]
		// costly but we don't want to lost values here
		if reflect.ValueOf(uk).IsZero() && !reflect.ValueOf(v).IsZero() {
			updates[k] = v
		}
	}

	return ds.Put(updates)
}

func (ds *DynamoService) Query(qps []QueryParams) ([]interface{}, error) {
	var expressions []queryExpression
	var attributes []map[string]interface{}

	for _, qp := range qps {
		var exist *TableSecondaryIndex = nil
		for _, ind := range ds.baseParams.indexes {
			if ind.key == qp.key {
				exist = &ind

				attributes = append(attributes, map[string]interface{}{
					":" + qp.key: qp.value,
				})

				expressions = append(expressions, queryExpression{
					IndexName:  exist.name,
					Expression: fmt.Sprintf("%s = :%s", qp.key, qp.key),
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
		TableName:                 &ds.baseParams.table,
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
		out, err := ds.svc.Query(&params)
		if ds.logger != nil {
			ds.logger.Println("QUERY: ", params)
		}

		if err != nil {
			return nil, err
		}

		if out.Items != nil {
			for _, outItem := range out.Items {
				newp := ds.empty()
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
