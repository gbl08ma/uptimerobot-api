package uptimerobot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// UptimeRobot is a representation of the UptimeRobot public API
type UptimeRobot struct {
	apikey     string
	HTTPClient *http.Client
}

// New creates a new UptimeRobot API client with the given API-key to identify
// the account you're working with
func New(apikey string) *UptimeRobot {
	return &UptimeRobot{
		apikey:     apikey,
		HTTPClient: http.DefaultClient,
	}
}

func (u *UptimeRobot) doRequest(apiMethod string, params *url.Values, target interface{}) error {
	if params == nil {
		params = &url.Values{}
	}

	params.Set("noJsonCallback", "1") // Enforce not to get JSONP wrapper
	params.Set("format", "json")
	params.Set("apiKey", u.apikey)

	url := url.URL{
		Scheme:   "https",
		Host:     "api.uptimerobot.com",
		Path:     fmt.Sprintf("/%s", apiMethod),
		RawQuery: params.Encode(),
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return err
	}

	res, err := u.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(target)
	return err
}

func (u *UptimeRobot) buildIntList(in interface{}) string {
	m := []string{}
	for _, i := range in.([]int) {
		m = append(m, strconv.FormatInt(int64(i), 10))
	}
	return strings.Join(m, "-")
}

func (u *UptimeRobot) bool2str(in bool) string {
	if in {
		return "1"
	}
	return "0"

}
