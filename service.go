package dynamo

import (
	"log"
	"os"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Config struct {
	Table            string
	Model            interface{}
	UseDescribeTable bool
}

type TableDef struct {
	baseParams tableSchema
	model      interface{}
	t          reflect.Type
}

type DynamoService struct{}

// Use global session to avoid creating new ones
var sess *session.Session = nil

var db *dynamodb.DynamoDB = nil

var logger *log.Logger = nil

func (d *DynamoService) Table(c Config) *TableDef {
	table := TableDef{}

	if c.UseDescribeTable {
		// Get description and at it to our service for later use
		table.baseParams = schemaFromDescribeTable(c.Table)
	} else {
		table.baseParams = schemaFromReflection(c.Table, c.Model)
	}

	table.t = reflect.TypeOf(c.Model)
	table.model = c.Model

	return &table
}

func New(s *session.Session, c *aws.Config) *DynamoService {
	if sess == nil {
		sess = s
	}

	if db == nil {
		db = dynamodb.New(sess, c)
	}

	if os.Getenv("DYNAMO_VERBOSE") != "" {
		logger = log.Default()
	}

	return &DynamoService{}
}
