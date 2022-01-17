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
	Table string
	Model interface{}
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
	return &TableDef{
		baseParams: schemaFromReflection(c.Table, c.Model),
		t:          reflect.TypeOf(c.Model),
		model:      c.Model,
	}
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
