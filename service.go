package dynamo

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type ObjectResolver interface {
	New() interface{}
	PrimaryKey(v interface{}) string
}

type Config struct {
	Table   string
	Log     bool
	Creater ObjectResolver
}

type DynamoService struct {
	svc        *dynamodb.DynamoDB
	baseParams TableDescription
	logger     *log.Logger
	creater    ObjectResolver
}

func New(c Config, options *aws.Config) *DynamoService {
	service := DynamoService{}
	var sess *session.Session
	if options != nil {
		sess = session.Must(session.NewSession(options))
	} else {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	}

	// get dynamo service
	svc := dynamodb.New(sess)
	service.svc = svc

	// get aws description
	description, err := svc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: &c.Table,
	})
	if err != nil {
		panic(err)
	}

	// Get description and at it to our service for later use
	service.baseParams = parseDescribeTable(description)
	service.creater = c.Creater
	if c.Log {
		service.logger = log.Default()
	}
	return &service
}
