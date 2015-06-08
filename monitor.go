package uptimerobot

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type MonitorStatus int

const (
	MonitorStatusPaused        MonitorStatus = 0
	MonitorStatusNotCheckedYet               = 1
	MonitorStatusUp                          = 2
	MonitorStatusSeemsDown                   = 8
	MonitorStatusDown                        = 9
)

type MonitorType int

const (
	_ MonitorType = iota
	MonitorTypeHTTP
	MonitorTypeKeyword
	MonitorTypePing
	MonitorTypePort
)

type MonitorSubtype int

const (
	_ MonitorSubtype = iota
	MonitorSubtypeHTTP
	MonitorSubtypeHTTPS
	MonitorSubtypeFTP
	MonitorSubtypeSMTP
	MonitorSubtypePOP3
	MonitorSubtypeIMAP
	MonitorSubtypeCustomPort = 99
)

type MonitorKeywordType int

const (
	_ MonitorKeywordType = iota
	MonitorKeywordTypeExists
	MonitorKeywordTypeNotExists
)

type Monitor struct {
	ID                 int                `json:"id,string"`
	FriendlyName       string             `json:"friendlyname"`
	URL                string             `json:"url"`
	Type               MonitorType        `json:"type,string"`
	Subtype            MonitorSubtype     `json:"subtype,string"`
	KeywordType        MonitorKeywordType `json:"keywordtype,string"`
	KeywordValue       string             `json:"keywordvalue"`
	HTTPUsername       string             `json:"httpusername"`
	HTTPPassword       string             `json:"httppassword"`
	Port               int                `json:"port,string"`
	Interval           int                `json:"interval,string"`
	Status             MonitorStatus      `json:"status,string"`
	AlltimeUptimeRatio float64            `json:"alltimeuptimeratio,string"`
	CustomUptimeRatio  float64            `json:"customuptimeratio,string"`
	AlertContacts      []AlertContact     `json:"alertcontact"`
	Logs               []Log              `json:"log"`
	ResponseTimes      []ResponseTime     `json:"responsetime"`
}

type ResponseTime struct {
	DateTime UptimeRobotDate `json:"datetime"`
	Value    int             `json:"value,string"`
}

// GetMonitorsInput is the Input type for the GetMonitors function. All parameters
// are optional and does not need to be set. The default values will give a list
// of all monitors
type GetMonitorsInput struct {
	// optional (if not used, will return all monitors in an account. Else, it is
	// possible to define any number of monitors with their IDs
	Monitors []int
	// optional (if not used, will return all monitors types (HTTP, keyword, ping..)
	// in an account. Else, it is possible to define any number of monitor types
	Types []MonitorType
	// optional (if not used, will return all monitors statuses (up, down, paused)
	// in an account. Else, it is possible to define any number of monitor statuses
	Statuses []MonitorStatus
	// optional (defines the number of days to calculate the uptime ratio(s) for
	// (Example: []int{7,30,45})
	CustomUptimeRatio []int
	// optional (defines if the logs of each monitor will be returned.)
	Logs bool
	// optional (defines if the response time data of each monitor will be returned.)
	ResponseTimes bool
	// optional (by default, response time value of each check is returned. The
	// API can return average values in given minutes. Default is 0. For ex: the
	// Uptime Robot dashboard displays the data averaged/grouped in 30 minutes)
	ResponseTimeAverage int
	// optional and works only for the Pro Plan as 24 hour+ logs are kept only in
	// the Pro Plan (starting date of the response times, must be used with
	// responseTimesEndDate) (can only be used if Monitors parameter is used with
	// a single monitorID and ResponseTimeEndDate - ResponseTimeStartDate can't
	// be more than 7 days)
	ResponseTimeStartDate *time.Time
	// optional and works only for the Pro Plan as 24 hour+ logs are kept only in
	// the Pro Plan (ending date of the response times, must be used with
	// responseTimesStartDate) (can only be used if Monitors parameter is used
	// with a single monitorID and ResponseTimeEndDate - ResponseTimeStartDate
	// can't be more than 7 days)
	ResponseTimeEndDate *time.Time
	// optional (defines if the notified alert contacts of each notification will
	// be returned.)
	LogAlertContacts bool
	// optional (defines if the alert contacts set for the monitor to be returned.)
	ShowMonitorAlertContacts bool
	// optional (defines if the user's timezone should be returned.)
	ShowTimezone bool
	// optional (a keyword of your choice to search within Monitor.URL and
	// Monitor.FriendlyName and get filtered results)
	Search string
}

// GetMonitors is a Swiss-Army knife type of a method for getting any information
// on monitors.
func (u *UptimeRobot) GetMonitors(in *GetMonitorsInput) ([]Monitor, error) {
	params := url.Values{}

	if len(in.Monitors) > 0 {
		params.Set("monitors", u.buildIntList(in.Monitors))
	}

	if len(in.Types) > 0 {
		params.Set("types", u.buildIntList(in.Types))
	}

	if len(in.Statuses) > 0 {
		params.Set("statuses", u.buildIntList(in.Statuses))
	}

	if len(in.CustomUptimeRatio) > 0 {
		params.Set("customUptimeRatio", u.buildIntList(in.CustomUptimeRatio))
	}

	params.Set("logs", u.bool2str(in.Logs))
	params.Set("responseTimes", u.bool2str(in.ResponseTimes))

	if in.ResponseTimeAverage > 0 {
		params.Set("responseTimeAverage", strconv.FormatInt(int64(in.ResponseTimeAverage), 10))
	}

	if in.ResponseTimeStartDate != nil && in.ResponseTimeEndDate != nil {
		if len(in.Monitors) != 1 || in.ResponseTimeEndDate.Sub(*in.ResponseTimeStartDate) > 7*24*time.Hour {
			return []Monitor{}, fmt.Errorf("Logic error. Please check documentation for StartDate & EndDate")
		}

		params.Set("responseTimesStartDate", in.ResponseTimeStartDate.Format("2006-01-02"))
		params.Set("responseTimesEndDate", in.ResponseTimeEndDate.Format("2006-01-02"))
	}

	params.Set("alertContacts", u.bool2str(in.LogAlertContacts))
	params.Set("showMonitorAlertContacts", u.bool2str(in.ShowMonitorAlertContacts))
	params.Set("showTimezone", u.bool2str(in.ShowTimezone))

	params.Set("offset", "0")
	params.Set("limit", "50")

	if in.Search != "" {
		params.Set("search", in.Search)
	}

	result := []Monitor{}
	for {
		res := &struct {
			Stat     string `json:"stat"`
			Offset   int    `json:"offset,string"`
			Limit    int    `json:"limit,string"`
			Total    int    `json:"total,string"`
			Monitors struct {
				Monitors []Monitor `json:"monitor"`
			} `json:"monitors"`
		}{}

		err := u.doRequest("getMonitors", &params, res)
		if err != nil {
			return []Monitor{}, err
		}

		if res.Stat != "ok" {
			return nil, fmt.Errorf("Got unexpected status: %s", res.Stat)
		}

		for _, m := range res.Monitors.Monitors {
			result = append(result, m)
		}

		if res.Offset+res.Limit >= res.Total {
			break
		}

		params.Set("offset", strconv.FormatInt(int64(res.Offset+res.Limit), 10))
	}

	return result, nil
}

// NewOrEditMonitor creates a new monitor if you do not pass an ID in the input,
// otherwise the monitor is updated
func (u *UptimeRobot) NewOrEditMonitor(in Monitor) (*Monitor, error) {
	params := &url.Values{}

	if in.FriendlyName == "" || in.URL == "" || in.Type == 0 {
		return nil, fmt.Errorf("Required parameters misisng. Please check the documentation.")
	}

	params.Set("monitorFriendlyName", in.FriendlyName)
	params.Set("monitorURL", in.URL)
	params.Set("monitorType", strconv.FormatInt(int64(in.Type), 10))

	if in.Subtype != 0 {
		params.Set("monitorSubType", strconv.FormatInt(int64(in.Subtype), 10))
	}

	if in.Port != 0 {
		params.Set("monitorPort", strconv.FormatInt(int64(in.Port), 10))
	}

	if in.KeywordType != 0 {
		params.Set("monitorKeywordType", strconv.FormatInt(int64(in.KeywordType), 10))
	}

	if in.KeywordValue != "" {
		params.Set("monitorKeywordValue", in.KeywordValue)
	}

	if in.HTTPUsername != "" {
		params.Set("monitorHTTPUsername", in.HTTPUsername)
	}

	if in.HTTPPassword != "" {
		params.Set("monitorHTTPPassword", in.HTTPPassword)
	}

	if len(in.AlertContacts) > 0 {
		m := []string{}
		for _, i := range in.AlertContacts {
			m = append(m, fmt.Sprintf("%d_%d_%d", i.ID, i.Threshold, i.Recurrence))
		}
		params.Set("monitorAlertContacts", strings.Join(m, "-"))
	}

	if in.Interval > 0 {
		params.Set("monitorInterval", strconv.FormatInt(int64(in.Interval), 10))
	}

	res := &struct {
		Stat    string  `json:"stat"`
		Monitor Monitor `json:"monitor"`
	}{
		Monitor: in,
	}

	var err error
	if in.ID == 0 {
		err = u.doRequest("newMonitor", params, res)
	} else {
		params.Set("monitorID", strconv.FormatInt(int64(in.ID), 10))
		err = u.doRequest("editMonitor", params, res)
	}
	if err != nil {
		return nil, err
	}

	if res.Stat == "ok" {
		return &res.Monitor, nil
	}

	return nil, fmt.Errorf("Got unexpected status: %s", res.Stat)
}

// DeleteMonitor deletes the monitor identifed by the monitorID
func (u *UptimeRobot) DeleteMonitor(monitorID int) error {
	res := &struct {
		Stat string `json:"stat"`
	}{}

	err := u.doRequest("deleteMonitor", &url.Values{
		"monitorID": []string{strconv.FormatInt(int64(monitorID), 10)},
	}, res)

	if err != nil {
		return err
	}

	if res.Stat == "ok" {
		return nil
	}

	return fmt.Errorf("Got unexpected status: %s", res.Stat)
}

// ResetMonitor will reset (deleting all stats and response time data) a monitor
func (u *UptimeRobot) ResetMonitor(monitorID int) error {
	res := &struct {
		Stat string `json:"stat"`
	}{}

	err := u.doRequest("resetMonitor", &url.Values{
		"monitorID": []string{strconv.FormatInt(int64(monitorID), 10)},
	}, res)

	if err != nil {
		return err
	}

	if res.Stat == "ok" {
		return nil
	}

	return fmt.Errorf("Got unexpected status: %s", res.Stat)
}
