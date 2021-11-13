package dynamo

import "github.com/aws/aws-sdk-go/service/dynamodb"

type TableSecondaryIndex struct {
	key  string
	name string
}

type TableDescription struct {
	arn     string
	table   string
	key     string
	indexes []TableSecondaryIndex
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

func parseDescribeTable(in *dynamodb.DescribeTableOutput) TableDescription {
	out := TableDescription{}

	out.table = *in.Table.TableName
	out.arn = *in.Table.TableArn

	// Get primary index
	out.key = hash(in.Table.KeySchema)

	// Get secondary indexes
	indexes := []TableSecondaryIndex{}
	for _, secondary := range in.Table.GlobalSecondaryIndexes {
		indexes = append(indexes, TableSecondaryIndex{
			key:  hash(secondary.KeySchema),
			name: *secondary.IndexName,
		})
	}

	out.indexes = indexes
	return out
}
