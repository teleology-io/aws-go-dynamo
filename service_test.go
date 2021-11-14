package dynamo

import (
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

func (inv Invitation) PrimaryKey(v interface{}) string {
	var invitation Invitation
	mapstructure.Decode(v, &invitation)
	return invitation.ID
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

var config = Config{
	Table:   "foundation-dev-invitations",
	Log:     false,
	Creater: Invitation{},
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

	// var out Invitation
	out, err := ddb.Get("bcf07e7e-c441-4eb6-8f94-8f1f14b369c5")
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
	ddb := New(config, nil)

	out := Invitation{
		ID:        TEST_INVITATION_ID,
		ProjectID: "085dea80-1e49-4cd1-a739-e9f355f08f2a",
		Role:      "viewer",
		Email:     "csullivan@teleology.io",
		Expiry:    time.Now().Unix(),
	}

	_, err := ddb.Put(out)
	if err != nil {
		t.Fail()
	}

	// We should have a collision here
	_, err = ddb.Create(out)
	if err == nil {
		t.Fail()
	}

	// Cleanup and test delete
	err = ddb.Delete(TEST_INVITATION_ID)
	if err != nil {
		t.Fail()
	}
}

func TestCreateAndUpdate(t *testing.T) {
	ddb := New(config, nil)

	_, err := ddb.Create(sample)
	if err != nil {
		t.Fail()
	}
	if err != nil {
		t.Fail()
	}

	updateResult, err := ddb.Update(Invitation{
		ID:   TEST_INVITATION_ID,
		Role: "publisher",
	})
	if err != nil {
		t.Fail()
	}

	// check we can convert to our expected type
	var invitation Invitation
	err = Convert(updateResult, &invitation)
	if err != nil {
		t.Fail()
	}
}

func TestQuery(t *testing.T) {
	ddb := New(config, nil)

	out, err := ddb.Query([]QueryParams{
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
		t.Fail()
	}

	invitations := make([]Invitation, len(out))
	err = Convert(&out, &invitations)
	if err != nil {
		t.Fail()
	}
}
