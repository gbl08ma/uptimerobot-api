package uptimerobot

import "time"

type UptimeRobotDate time.Time

func (t UptimeRobotDate) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).Format("\"01/02/2006 15:04:05\"")), nil
}

func (t *UptimeRobotDate) UnmarshalJSON(in []byte) error {
	p, err := time.Parse("\"01/02/2006 15:04:05\"", string(in))
	if err != nil {
		// Some API parts are delivering a different date format so we handle this too
		p, err = time.Parse("\"01/02/06 15:04:05\"", string(in))
		if err != nil {
			return err
		}
	}
	*t = UptimeRobotDate(p)
	return nil
}

func (t UptimeRobotDate) String() string {
	return time.Time(t).String()
}
