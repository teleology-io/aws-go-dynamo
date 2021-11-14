package dynamo

import (
	"fmt"
	"testing"
	"time"

	"github.com/mitchellh/mapstructure"
)

var TEST_INVITATION_ID = "test-invitation-id"

type Invitation struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Expiry    int64  `json:"expiry"`
}

func (v Invitation) New() interface{} {
	return Invitation{}
}

var config = Config{
	table:   "foundation-dev-invitations",
	log:     false,
	creater: Invitation{},
}

var sample = Invitation{
	ID:        TEST_INVITATION_ID,
	ProjectID: "085dea80-1e49-4cd1-a739-e9f355f08f2a",
	Role:      "viewer",
	Email:     "csullivan@teleology.io",
	Expiry:    time.Now().Unix(),
}

func Convert(i interface{}, o interface{}) error {
	return mapstructure.Decode(i, o)
}

func TestGet(t *testing.T) {
	ddb := New(config, nil)

	// var out Invitation
	out, err := ddb.Get("bcf07e7e-c441-4eb6-8f94-8f1f14b369c5")
	if err != nil {
		t.Fail()
	}

	var invitation Invitation
	err = Convert(out, &invitation)
	if err != nil {
		t.Fail()
	}
	fmt.Println("out", invitation)
}

func TestCollisions(t *testing.T) {
	ddb := New(config, nil)

	out := Invitation{
		ID:        TEST_INVITATION_ID,
		ProjectID: "085dea80-1e49-4cd1-a739-e9f355f08f2a",
		Role:      "viewer",
		Email:     "csullivan@teleology.io",
		Expiry:    time.Now().Unix(),
	}

	putResult, err := ddb.Put(out)
	if err != nil {
		t.Fail()
	}

	// We should have a collision here
	createResult, err := ddb.Create(TEST_INVITATION_ID, out)
	if err == nil {
		t.Fail()
	}

	// Cleanup and test delete
	deleteResult, err := ddb.Delete(TEST_INVITATION_ID)
	if err != nil {
		t.Fail()
	}

	fmt.Println("UPDATE", putResult)
	fmt.Println("CREATE", createResult)
	fmt.Println("DELETE", deleteResult)
}

func TestCreate(t *testing.T) {
	ddb := New(config, nil)

	createResult, err := ddb.Create(TEST_INVITATION_ID, sample)
	if err != nil {
		t.Fail()
	}
	if err != nil {
		t.Fail()
	}

	newer := Invitation{
		ID:   TEST_INVITATION_ID,
		Role: "publisher",
	}

	updateResult, err := ddb.Update(TEST_INVITATION_ID, newer)
	if err != nil {
		t.Fail()
	}

	fmt.Println("CREATE", createResult)
	fmt.Println("UPDATE", updateResult)
}

func TestQuery(t *testing.T) {
	ddb := New(config, nil)

	out, err := ddb.Query([]QueryParams{
		{
			key:   "project_id",
			value: "085dea80-1e49-4cd1-a739-e9f355f08f2a",
		},
		{
			key:   "email",
			value: "csullivan@teleology.io",
		},
	})
	if err != nil {
		t.Fail()
	}

	invitations := make([]Invitation, len(out))
	err = Convert(&out, &invitations)
	if err != nil {
		t.Fail()
	}

	fmt.Println("QUERY", invitations)
}
