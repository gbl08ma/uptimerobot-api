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
	ID           int                `json:"id,string"`
	Type         AlertContactType   `json:"type,string"`
	Value        string             `json:"value"`
	Status       AlertContactStatus `json:"status,string"`
	Threshold    int                `json:"threshold,string"`
	Recurrence   int                `json:"recurrence,string"`
	FriendlyName string             `json:"friendlyname"`
}

// GetAlertContacts can be used to retrieve a (filtered) list of alert contacts
func (u *UptimeRobot) GetAlertContacts(contactIDs []int) ([]AlertContact, error) {
	params := &url.Values{
		"limit":  []string{"50"},
		"offset": []string{"0"},
	}

	if len(contactIDs) > 0 {
		params.Set("alertcontacts", u.buildIntList(contactIDs))
	}

	response := []AlertContact{}

	for {
		res := &struct {
			Stat          string `json:"stat"`
			Offset        int    `json:"offset,string"`
			Limit         int    `json:"limit,string"`
			Total         int    `json:"total,string"`
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
		Message      string       `json:"message"`
		AlertContact AlertContact `json:"alertcontact"`
	}{
		AlertContact: in,
	}

	if in.Type == 0 || in.Value == "" {
		return nil, fmt.Errorf("Required parameters misisng. Please check the documentation.")
	}

	if len(in.FriendlyName) > 30 {
		return nil, fmt.Errorf("FriendlyName may not have more than 30 chars")
	}

	params.Set("alertContactType", strconv.FormatInt(int64(in.Type), 10))
	params.Set("alertContactValue", in.Value)

	if in.FriendlyName != "" {
		params.Set("alertContactFriendlyName", in.FriendlyName)
	}

	err := u.doRequest("newAlertContact", params, res)
	if err != nil {
		return nil, err
	}

	if res.Stat != "ok" {
		return nil, fmt.Errorf("Contact not created: %s", res.Message)
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
