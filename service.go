package dynamo

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type ObjectCreator interface {
	New() interface{}
}

type Config struct {
	table   string
	log     bool
	creater ObjectCreator
}

type DynamoService struct {
	svc        *dynamodb.DynamoDB
	baseParams TableDescription
	logger     *log.Logger
	creater    ObjectCreator
}

func (d DynamoService) New() interface{} {
	var data interface{}
	return data
}

func New(c Config, options *aws.Config) *DynamoService {
	service := DynamoService{}

	var sess *session.Session
	if options != nil {
		opts := session.Options{}
		opts.Config.MergeIn(options)
		sess = session.Must(session.NewSessionWithOptions(opts))
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
		TableName: &c.table,
	})
	if err != nil {
		panic(err)
	}

	// Get description and at it to our service for later use
	service.baseParams = parseDescribeTable(description)
	service.creater = c.creater
	if c.log {
		service.logger = log.Default()
	}
	return &service
}
