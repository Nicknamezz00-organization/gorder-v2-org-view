package consts

const (
	ErrnoSuccess      = 0
	ErrnoUnknownError = 1

	// param error 1xxx
	ErrnoBindRequestError     = 1000
	ErrnoRequestValidateError = 1001

	// mysql error 2xxx
)

var ErrMsg = map[int]string{
	ErrnoSuccess:      "success",
	ErrnoUnknownError: "unknown error",

	ErrnoBindRequestError:     "bind request error",
	ErrnoRequestValidateError: "validate request error",
}
