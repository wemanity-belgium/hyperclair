package xerrors

import "errors"

var (
	ServerUnavailable = errors.New("Cannot reach Hyperclair server. Is it running?")
	Unauthorized      = errors.New("Unauthorized access")
	InternalError     = errors.New("Client quit unexpectedly")
	NotFound          = errors.New("Image not found")
)
