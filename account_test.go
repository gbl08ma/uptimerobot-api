package uptimerobot

import (
	"strings"
	"testing"
)

func TestGetAccountDetail(t *testing.T) {
	ur := New("u232958-fc43e2ab62ed66a08b0e578b")
	ur.FullDebug = true
	ad, err := ur.GetAccountDetails()

	if err != nil {
		t.Fatalf("Test errored: %s", err)
	}

	if ad.MonitorLimit < 50 {
		t.Errorf("MonitorLimit should be >= 50, is %d", ad.MonitorLimit)
	}

	if ad.MonitorInterval > 5 {
		t.Errorf("MonitorInterval should be <= 5, is %d", ad.MonitorInterval)
	}
}

func TestGetAccountDetailWithoutAccount(t *testing.T) {
	ur := New("foobar")
	ur.FullDebug = true
	_, err := ur.GetAccountDetails()

	if err == nil {
		t.Errorf("Test should have errored.")
	}

	if !strings.HasPrefix(err.Error(), "Got unexpected status:") {
		t.Errorf("Test errored with unexpected result: %s", err)
	}
}
