package uptimerobot

type APIError int

const (
	ErrorAPIKeyWrongFormat APIError = 100 // apiKey not mentioned or in a wrong format
	ErrorAPIKeyWrong                = 101 // apiKey is wrong
	ErrorWrongFormat                = 102 // format is wrong (should be xml or json)
	ErrorNoSuchMethod               = 103 // No such method exists

	ErrorMonitorIDShouldBeInteger                  = 200 // monitorID(s) should be integers
	ErrorMonitorURLInvalid                         = 201 // monitorUrl is invalid
	ErrorMonitorTypeInvalid                        = 202 // monitorType is invalid
	ErrorMonitorSubTypeInvalid                     = 203 // monitorSubType is invalid
	ErrorMonitorKeywordTypeInvalid                 = 204 // monitorKeywordType is invalid
	ErrorMonitorPortInvalid                        = 205 // monitorPort is invalid
	ErrorMonitorFriendlyNameRequired               = 206 // monitorFriendlyName is required
	ErrorMonitorAlreadyExists                      = 207 // The monitor already exists
	ErrorMonitorSubTypeRequired                    = 208 // monitorSubType is required for this type of monitors
	ErrorMonitorKeywordTypeAndKeywordValueRequired = 209 // monitorKeyWordType and monitorKeyWordValue are required for this type of monitors
	ErrorMonitorIDNoExists                         = 210 // monitorID doesn't exist
	ErrorMonitorIDRequired                         = 211 // monitorID is required
	ErrorAccountHasNoMonitors                      = 212 // The account has no monitors
	ErrorNoEditsFound                              = 213 // At least one of the parameters to be edited are required
	ErrorHTTPCredentialsMismatch                   = 214 // monitorHTTPUsername and monitorHTTPPassword should both be empty or have values
	ErrorInvalidAPIScope                           = 215 // monitor specific apiKeys can only use getMonitors method
	ErrorEMailInUse                                = 216 // A user with this e-mail already exists
	ErrorFirstLastNameEMailRequired                = 217 // userFirstLastName and userEmail are both required
	ErrorEMailFormatInvalid                        = 218 // userEmail is not in the right e-mail format
	ErrorUserCreateNotAllowed                      = 219 // This account is not authorized to create users
	ErrorMonitorAlertContactsValueInvalid          = 220 // monitorAlertContacts value is wrong
	ErrorNoAlertContactsFound                      = 221 // The account has no alert contacts
	ErrorAlertContactIDShoudBeInteger              = 222 // alertcontactID(s) should be integers
	ErrorAlertContactTypeAndValueRequired          = 223 // alertContactType and alertContactValue are both required
	ErrorAlertContactTypeNotSupported              = 224 // This alertContactType is not supported"
	ErrorAlertContactAlreadyExists                 = 225 // The alert contact already exists
	ErrorAlertContactDoesNotFollowUptimeRobot      = 226 // The alert contact is not following @uptimerobot Twitter user. It is required so that the Twitter direct messages (DM) can be sent
	ErrorBoxcarUserNotExists                       = 227 // The Boxcar user mentioned does not exist
	ErrorBoxcarUserNotAdded                        = 228 // The Boxcar alert contact couldn't be added, please try again later
	ErrorAlertContactIDNotExists                   = 229 // alertContactID doesn't exist
	ErrorAlertContactValueShouldBeEMail            = 230 // alertContactValue should be a valid e-mail for this alertContactType
)
