package constants

type ctxKeyRequestID int

const RequestIDKey ctxKeyRequestID = 0

type ctxKey string

const CountryCodeKey ctxKey = "country_code"

const (
	RequestIDHeader = "X-Request-Id"
)

const (
	SendOTP = "send_otp"
)

const (
	AccessSubject  = "access"
	RefreshSubject = "refresh"
)

const (
	SSLModeDisable = "disable"
)
