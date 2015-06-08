package uptimerobot

import "fmt"

// AccountDetail represents detailed information about the account
type AccountDetail struct {
	// the max number of monitors that can be created for the account
	MonitorLimit int `json:"monitorLimit,string"`
	// the min monitoring interval (in minutes) supported by the account
	MonitorInterval int `json:"monitorInterval,string"`
	// the number of "up" monitors
	UpMonitors int `json:"upMonitors,string"`
	// the number of "down" monitors
	DownMonitors int `json:"downMonitors,string"`
	// the number of "paused" monitors
	PausedMonitors int `json:"pausedMonitors,string"`
}

// GetAccountDetails fetches account details (max number of monitors that can be
// added and number of up/down/paused monitors) about the account identified by
// the given APIKey
func (u *UptimeRobot) GetAccountDetails() (*AccountDetail, error) {
	result := &struct {
		Stat    string        `json:"stat"`
		Account AccountDetail `json:"account"`
	}{}

	err := u.doRequest("getAccountDetails", nil, result)
	if err != nil {
		return nil, err
	}

	if result.Stat == "ok" {
		return &result.Account, nil
	}

	return nil, fmt.Errorf("Got unexpected status: %s", result.Stat)
}
