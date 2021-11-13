package dynamo

import (
	"fmt"
	"testing"
	"time"
)

var TEST_INVITATION_ID = "test-invitation-id"

type Invitation struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Expiry    int64  `json:"expiry"`
}

var config = Config{
	table: "foundation-dev-invitations",
	log:   false,
	empty: func() interface{} {
		return Invitation{}
	},
}

var sample = Invitation{
	ID:        TEST_INVITATION_ID,
	ProjectID: "085dea80-1e49-4cd1-a739-e9f355f08f2a",
	Role:      "viewer",
	Email:     "csullivan@teleology.io",
	Expiry:    time.Now().Unix(),
}

func TestGet(t *testing.T) {
	ddb := New(config, nil)

	var out Invitation
	err := ddb.Get("bcf07e7e-c441-4eb6-8f94-8f1f14b369c5", &out)
	if err != nil {
		t.Fail()
	}
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

	err := ddb.Put(out)
	if err != nil {
		t.Fail()
	}

	// We should have a collision here
	err = ddb.Create(TEST_INVITATION_ID, out)
	if err == nil {
		t.Fail()
	}

	// Cleanup and test delete
	_, err = ddb.Delete(TEST_INVITATION_ID)
	if err != nil {
		t.Fail()
	}
}

func TestCreate(t *testing.T) {
	ddb := New(config, nil)

	err := ddb.Create(TEST_INVITATION_ID, sample)
	if err != nil {
		t.Fail()
	}

	newer := Invitation{
		ID:   TEST_INVITATION_ID,
		Role: "publisher",
	}

	ddb.Update(TEST_INVITATION_ID, newer)
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

	// casted := Invitation(out[0].(Invitation))

	fmt.Println("casted", out)

	if err != nil {
		t.Fail()
	}
}
