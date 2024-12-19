package constants

// User
const (
	ErrMsgInvalidCredentials = "invalid email or password"
	ErrMsgUserNotFound       = "user not found"
	ErrMsgEmailInUse         = "email is already in use"
	ErrMsgCPFInUse           = "CPF is already in use"
	ErrMsgInvalidEmail       = "invalid email"
	ErrMsgInvalidCPF         = "invalid CPF"
	ErrMsgInvalidAge         = "invalid age: must be between 0 and 150"
	ErrMsgInvalidName        = "invalid name: must be at least 2 characters"
)

// Auth
const (
	ErrMsgMissingHeader    = "missing authorization header"
	ErrMsgInvalidToken     = "invalid or expired token"
	ErrMsgInvalidTokenType = "invalid token type"
)

const PORT = ":3000"

const TRACER_NAME = "golerplate"
