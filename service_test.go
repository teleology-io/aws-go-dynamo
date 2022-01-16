package dynamo

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/mitchellh/mapstructure"
)

var TEST_INVITATION_ID = "test-invitation-id"

type Invitation struct {
	ID        string `json:"id" dynamo:"pk"`
	ProjectID string `json:"project_id" dynamo:"gsi,ProjectGSI"`
	Role      string `json:"role"`
	Email     string `json:"email" dynamo:"gsi,EmailGSI"`
	Expiry    int64  `json:"expiry"`
}

func Convert(in interface{}, out interface{}) error {
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

var sample = Invitation{
	ID:        TEST_INVITATION_ID,
	ProjectID: "085dea80-1e49-4cd1-a739-e9f355f08f2a",
	Role:      "viewer",
	Email:     "csullivan@teleology.io",
	Expiry:    time.Now().Unix(),
}

var table *tableDef = nil

func init() {
	var endpoint = "http://localhost:8000"
	var region = "us-east-1"

	sess := session.Must(session.NewSession())

	db := New(sess, &aws.Config{
		Endpoint: &endpoint,
		Region:   &region,
	})

	table = db.Table(Config{
		Table:            "foundation-local-invitations",
		Model:            Invitation{},
		UseDescribeTable: true,
	})

	_, err := table.Put(sample)
	if err != nil {
		fmt.Println("err", err.Error())
		panic("could not initialize tests")
	}
}

func TestGet(t *testing.T) {
	// var out Invitation
	out, err := table.Get(TEST_INVITATION_ID)
	if err != nil {
		t.Fail()
	}

	// check we can convert to our expected type
	var invitation Invitation
	err = Convert(out, &invitation)
	if err != nil {
		t.Fail()
	}
}

func TestCollisions(t *testing.T) {
	out := Invitation{
		ID:        TEST_INVITATION_ID,
		ProjectID: "085dea80-1e49-4cd1-a739-e9f355f08f2a",
		Role:      "viewer",
		Email:     "csullivan@teleology.io",
		Expiry:    time.Now().Unix(),
	}

	_, err := table.Put(out)
	if err != nil {
		t.Fail()
	}

	// We should have a collision here
	_, err = table.Create(out.ID, out)
	if err == nil {
		t.Fail()
	}

	// Cleanup and test delete
	err = table.Delete(TEST_INVITATION_ID)
	if err != nil {
		t.Fail()
	}
}

func TestCreateAndUpdate(t *testing.T) {
	_, err := table.Create(sample.ID, sample)
	if err != nil {
		fmt.Println("creat er", err.Error())
		t.Fail()
	}

	updateResult, err := table.Update(TEST_INVITATION_ID, Invitation{
		ID:   TEST_INVITATION_ID,
		Role: "publisher",
	})
	if err != nil {
		fmt.Println("er", err.Error())
		t.Fail()
	}

	// check we can convert to our expected type
	invitation := Invitation{}
	err = Convert(updateResult, &invitation)
	if err != nil {
		t.Fail()
	}
}

func TestQuery(t *testing.T) {
	out, err := table.Query([]QueryParams{
		{
			Key:   "project_id",
			Value: "085dea80-1e49-4cd1-a739-e9f355f08f2a",
		},
		{
			Key:   "email",
			Value: "csullivan@teleology.io",
		},
	})
	if err != nil {
		fmt.Println("query out", err.Error())
		t.Fail()
	}

	invitations := make([]Invitation, len(out))
	err = Convert(&out, &invitations)
	if err != nil {
		fmt.Println("query convert", err.Error())
		t.Fail()
	}
}
