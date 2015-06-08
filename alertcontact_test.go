package uptimerobot

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/satori/go.uuid"
)

func TestGetAlertContacts(t *testing.T) {
	ur := New(os.Getenv("UR_API_KEY"))
	ur.FullDebug = true

	ac, err := ur.GetAlertContacts(nil)
	if err != nil {
		t.Fatalf("Test errored: %s", err)
	}

	if len(ac) == 0 {
		t.Errorf("Initial count of contacts is expected to be > 0, is %d", len(ac))
	}

	var mailContact *AlertContact
	for _, c := range ac {
		if c.Type == AlertContactTypeEMail {
			mailContact = &c
		}
	}

	if mailContact == nil {
		t.Errorf("Account is expected to have at least one mail contact")
	}
}

func TestNewGetDeleteAlertContact(t *testing.T) {
	ur := New(os.Getenv("UR_API_KEY"))
	ur.FullDebug = true

	u := uuid.NewV4().String()
	ac := AlertContact{
		Type:         AlertContactTypeEMail,
		Value:        fmt.Sprintf("%s@example.com", u),
		FriendlyName: u[0:30],
	}

	ac2, err := ur.NewAlertContact(ac)
	if err != nil {
		t.Fatalf("Test errored: %s", err)
	}

	if ac2.ID == 0 {
		t.Fatalf("Expected to get a contact ID, got 0")
	}

	if ac2.FriendlyName != ac.FriendlyName || ac2.Type != ac.Type || ac2.Value != ac.Value {
		t.Errorf("Out contact did not match in contact: %+v != %+v", ac, ac2)
	}

	tmp, err := ur.GetAlertContacts([]int{ac2.ID})
	if err != nil {
		t.Fatalf("Test errored: %s", err)
	}

	if len(tmp) != 1 {
		t.Errorf("Expected to get one contact by ID, got %d", len(tmp))
	}

	err = ur.DeleteAlertContact(ac2.ID)
	if err != nil {
		t.Fatalf("Deleting the contact errored: %s", err)
	}

	tmp, err = ur.GetAlertContacts([]int{ac2.ID})
	if err != nil {
		t.Fatalf("Test errored: %s", err)
	}

	if len(tmp) != 0 {
		t.Errorf("Found alert contact after delete")
	}
}

func TestNewAlertContactMissingParameters(t *testing.T) {
	ur := New(os.Getenv("UR_API_KEY"))
	ur.FullDebug = true

	ac := AlertContact{}
	_, err := ur.NewAlertContact(ac)
	if err == nil {
		t.Fatalf("Test should have errored.")
	}

	if err.Error() != "Required parameters misisng. Please check the documentation." {
		t.Errorf("Got an unexpected error string: %s", err)
	}
}

func TestNewAlertContactWrongParameters(t *testing.T) {
	ur := New(os.Getenv("UR_API_KEY"))
	ur.FullDebug = true

	ac := AlertContact{
		Type:  -1,
		Value: "-1",
	}

	_, err := ur.NewAlertContact(ac)
	if err == nil {
		t.Fatalf("Test should have errored.")
	}

	if !strings.HasPrefix(err.Error(), "Contact not created") {
		t.Errorf("Got an unexpected error string: %s", err)
	}
}

func TestNewAlertContactLongFriendlyName(t *testing.T) {
	ur := New(os.Getenv("UR_API_KEY"))
	ur.FullDebug = true

	ac := AlertContact{
		FriendlyName: "I'm a very long friendly name with more than 30 characters which should fail",
		Value:        "foo",
		Type:         AlertContactTypeBoxcar,
	}

	_, err := ur.NewAlertContact(ac)
	if err == nil {
		t.Fatalf("Test should have errored.")
	}

	if err.Error() != "FriendlyName may not have more than 30 chars" {
		t.Errorf("Got an unexpected error string: %s", err)
	}
}
