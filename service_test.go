package dynamo

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var TEST_INVITATION_ID = "test-invitation-id"

type Invitation struct {
	ID        string `json:"id" dynamo:"pk"`
	ProjectID string `json:"project_id" dynamo:"gsi,ProjectGSI"`
	Role      string `json:"role"`
	Email     string `json:"email" dynamo:"gsi,EmailGSI"`
	Expiry    int64  `json:"expiry"`
}

var sample = Invitation{
	ID:        TEST_INVITATION_ID,
	ProjectID: "085dea80-1e49-4cd1-a739-e9f355f08f2a",
	Role:      "viewer",
	Email:     "csullivan@teleology.io",
	Expiry:    time.Now().Unix(),
}

var table *TableDef = nil

func init() {
	var endpoint = "http://localhost:8000"
	var region = "us-east-1"

	sess := session.Must(session.NewSession())

	db := New(sess, &aws.Config{
		Endpoint: &endpoint,
		Region:   &region,
	})

	table = db.Table(Config{
		Table: "foundation-local-invitations",
		Model: Invitation{},
	})

	_, err := table.Put(sample)
	if err != nil {
		fmt.Println("err", err.Error())
		panic("could not initialize tests")
	}
}

func TestGet(t *testing.T) {
	// var out Invitation
	_, err := table.Get(TEST_INVITATION_ID)
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
	_, err = table.Create(out)
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
	_, err := table.Create(sample)
	if err != nil {
		t.Fail()
	}

	_, err = table.Update(Invitation{
		ID:   TEST_INVITATION_ID,
		Role: "publisher",
	})
	if err != nil {
		t.Fail()
	}
}

func TestQuery(t *testing.T) {
	_, err := table.Query([]QueryParams{
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
}
