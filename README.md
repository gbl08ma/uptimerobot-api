# Jimdo / uptimerobot-api

[![Build Status](https://magnum.travis-ci.com/Jimdo/uptimerobot-api.svg?token=oqHdLv3md1svjAXzqgy2&branch=master)](https://magnum.travis-ci.com/Jimdo/uptimerobot-api)
[![API-Documentation](http://badge.luzifer.io/v1/badge?title=API&text=Documentation&color=4c1)](https://uptimerobot.com/api)
[![GoDoc Reference](http://badge.luzifer.io/v1/badge?color=5d79b5&title=godoc&text=reference)](https://godoc.org/github.com/Jimdo/uptimerobot-api)
[![License Apache 2.0](http://badge.luzifer.io/v1/badge?color=5d79b5&title=license&text=Apache%202.0)](http://www.apache.org/licenses/LICENSE-2.0)

UptimeRobot is an easy-to-use monitoring service. This library enables developers to use Go to access the API of UptimeRobot to manage their resources.

## Testing

To execute the tests you need to export your UptimeRobot API-Key to your env before running the tests:

```bash
# export UR_API_KEY=u232958-fc43e2ab62ed66a08b0e578b
# go test -cover -v ./... 2>/dev/null
=== RUN TestGetAccountDetail
--- PASS: TestGetAccountDetail (0.61s)
=== RUN TestGetAccountDetailWithoutAccount
--- PASS: TestGetAccountDetailWithoutAccount (0.54s)
=== RUN TestGetAlertContacts
--- PASS: TestGetAlertContacts (0.56s)
=== RUN TestNewGetDeleteAlertContact
--- PASS: TestNewGetDeleteAlertContact (2.23s)
=== RUN TestNewAlertContactMissingParameters
--- PASS: TestNewAlertContactMissingParameters (0.00s)
=== RUN TestNewAlertContactWrongParameters
--- PASS: TestNewAlertContactWrongParameters (0.54s)
=== RUN TestNewAlertContactLongFriendlyName
--- PASS: TestNewAlertContactLongFriendlyName (0.00s)
=== RUN TestMonitorFlow
--- PASS: TestMonitorFlow (3.89s)
PASS
coverage: 74.0% of statements
ok  	github.com/Jimdo/uptimerobot-api	8.385s
```
