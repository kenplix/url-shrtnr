package errorcode

type ErrorCode string

const (
	ParsingError         ErrorCode = "PARSING_ERROR"
	InvalidSchema        ErrorCode = "INVALID_SCHEMA"
	InvalidField         ErrorCode = "INVALID_FIELD"
	MissingField         ErrorCode = "MISSING_FIELD"
	AlreadyExists        ErrorCode = "ALREADY_EXISTS"
	IncorrectCredentials ErrorCode = "INCORRECT_CREDENTIALS"
	UnauthorizedAccess   ErrorCode = "UNAUTHORIZED_ACCESS"
	InternalError        ErrorCode = "INTERNAL_ERROR"
)
