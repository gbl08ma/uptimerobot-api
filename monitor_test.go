package uptimerobot

import (
	"fmt"
	"os"
	"testing"

	"github.com/satori/go.uuid"
)

func TestMonitorFlow(t *testing.T) {
	ur := New(os.Getenv("UR_API_KEY"))
	ur.FullDebug = true

	ac := setUpAlertContact(t, ur)
	monitor := createNewMonitor(t, ur, ac)
	findAndMatchMonitor(t, ur, monitor, ac)
	monitor = updateMonitor(t, ur, monitor)
	findAndMatchMonitor(t, ur, monitor, ac)
	deleteMonitor(t, ur, monitor)
	deleteAlertContact(t, ur, ac)
}

func setUpAlertContact(t *testing.T, ur *UptimeRobot) *AlertContact {
	ac := AlertContact{
		Type:  AlertContactTypeEMail,
		Value: fmt.Sprintf("%s@example.com", uuid.NewV4().String()),
	}

	ac2, err := ur.NewAlertContact(ac)
	if err != nil {
		t.Fatalf("Unable to create required contact: %s", err)
	}

	return ac2
}

func createNewMonitor(t *testing.T, ur *UptimeRobot, ac *AlertContact) *Monitor {
	monitor := Monitor{
		FriendlyName:  uuid.NewV4().String()[0:30],
		URL:           "http://www.example.com/",
		Type:          MonitorTypeHTTP,
		KeywordType:   MonitorKeywordTypeNotExists,
		KeywordValue:  "Example Domain",
		Interval:      10,
		AlertContacts: []AlertContact{*ac},
	}

	outMonitor, err := ur.NewOrEditMonitor(monitor)
	if err != nil {
		t.Fatalf("Test errored: %s", err)
	}

	if outMonitor.FriendlyName != monitor.FriendlyName || outMonitor.URL != monitor.URL || outMonitor.Type != monitor.Type || outMonitor.KeywordValue != monitor.KeywordValue {
		t.Errorf("Output monitor did not match input monitor: %+v != %+v", monitor, outMonitor)
	}

	if outMonitor.ID == 0 {
		t.Errorf("Expected to get a monitor ID but got none")
	}

	return outMonitor
}

func findAndMatchMonitor(t *testing.T, ur *UptimeRobot, mon *Monitor, ac *AlertContact) {
	mons, err := ur.GetMonitors(&GetMonitorsInput{
		ShowMonitorAlertContacts: true,
		ResponseTimes:            true,
	})

	if err != nil {
		t.Fatalf("Test errored: %s", err)
	}

	var outMonitor *Monitor
	for _, m := range mons {
		if m.ID == mon.ID {
			outMonitor = &m
			break
		}
	}

	if outMonitor == nil {
		t.Fatalf("Did not find specified monitor")
	}

	if len(outMonitor.AlertContacts) != 1 {
		t.Errorf("Expected to have 1 alert contact for this service, gpt %d", len(outMonitor.AlertContacts))
	}

	if outMonitor.AlertContacts[0].ID != ac.ID {
		t.Errorf("Did not have expected alert contact assigned to monitor.")
	}

	if outMonitor.FriendlyName != mon.FriendlyName {
		t.Errorf("Input did not match output: %+v != %+v", mon, outMonitor)
	}
}

func updateMonitor(t *testing.T, ur *UptimeRobot, mon *Monitor) *Monitor {
	mon.FriendlyName = uuid.NewV4().String()[0:30]

	m, err := ur.NewOrEditMonitor(*mon)
	if err != nil {
		t.Fatalf("Test errored: %s", err)
	}

	if m.ID != mon.ID {
		t.Errorf("A new monitor was created instead of an update")
	}

	return m
}

func deleteMonitor(t *testing.T, ur *UptimeRobot, mon *Monitor) {
	err := ur.DeleteMonitor(mon.ID)
	if err != nil {
		t.Fatalf("Test errored: %s", err)
	}
}

func deleteAlertContact(t *testing.T, ur *UptimeRobot, ac *AlertContact) {
	ur.DeleteAlertContact(ac.ID)
}
