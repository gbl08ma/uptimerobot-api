package uptimerobot

import (
	"fmt"
	"net/url"
	"strconv"
)

// AlertContactType represents the type of the alert contact
type AlertContactType int

const (
	_ AlertContactType = iota
	AlertContactTypeSMS
	AlertContactTypeEMail
	AlertContactTypeTwitterDM
	AlertContactTypeBoxcar
	AlertContactTypeWebHook
	AlertContactTypePushBullet
	AlertContactTypeZapier
	AlertContactTypePushover
	AlertContactTypeHipChat
	AlertContactTypeSlack
)

type AlertContactStatus int

const (
	AlertContactStatusNotActivated AlertContactStatus = iota
	AlertContactStatusPaused
	AlertContactStatusActive
)

type AlertContact struct {
	ID           int
	Type         AlertContactType
	Value        string
	Status       AlertContactStatus
	Threshold    int
	Recurrence   int
	FriendlyName string
}

// GetAlertContacts can be used to retrieve a (filtered) list of alert contacts
func (u *UptimeRobot) GetAlertContacts(contactIDs *[]int) ([]AlertContact, error) {
	params := &url.Values{
		"limit":  []string{"50"},
		"offset": []string{"0"},
	}

	if contactIDs != nil && len(*contactIDs) > 0 {
		params.Set("alertcontacts", u.buildIntList(contactIDs))
	}

	response := []AlertContact{}

	for {
		res := &struct {
			Stat          string `json:"stat"`
			Offset        int    `json:"offset"`
			Limit         int    `json:"limit"`
			Total         int    `json:"total"`
			AlertContacts struct {
				Contacts []AlertContact `json:"alertcontact"`
			} `json:"alertcontacts"`
		}{}

		err := u.doRequest("getAlertContacts", params, res)
		if err != nil {
			return []AlertContact{}, err
		}

		for _, i := range res.AlertContacts.Contacts {
			response = append(response, i)
		}

		if res.Limit+res.Offset >= res.Total {
			break
		}

		params.Set("offset", strconv.FormatInt(int64(res.Offset+res.Limit), 10))
	}

	return response, nil
}

// NewAlertContact creates a new alert contact of any type (mobile/SMS alert
// contacts are not supported yet)
func (u *UptimeRobot) NewAlertContact(in AlertContact) (*AlertContact, error) {
	params := &url.Values{}
	res := &struct {
		Stat         string       `json:"stat"`
		AlertContact AlertContact `json:"alertcontact"`
	}{
		AlertContact: in,
	}

	if in.Type == 0 || in.Value == "" {
		return nil, fmt.Errorf("Required parameters misisng. Please check the documentation.")
	}

	params.Set("alertContactType", strconv.FormatInt(int64(in.Type), 10))
	params.Set("alertContactValue", in.Value)

	if in.FriendlyName != "" {
		params.Set("alertContactFriendleName", in.FriendlyName)
	}

	err := u.doRequest("newAlertContact", params, res)
	if err != nil {
		return nil, err
	}

	return &res.AlertContact, nil
}

// DeleteAlertContact can be used to delete an alert contact
func (u *UptimeRobot) DeleteAlertContact(contactID int) error {
	res := &struct {
		Stat string `json:"stat"`
	}{}

	err := u.doRequest("deleteAlertContact", &url.Values{
		"alertContactID": []string{strconv.FormatInt(int64(contactID), 10)},
	}, res)

	if err != nil {
		return err
	}

	if res.Stat == "ok" {
		return nil
	}

	return fmt.Errorf("Got unexpected status: %s", res.Stat)
}
