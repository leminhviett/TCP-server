package errors

import "errors"

var (
	ErrorDialConn = errors.New("Error when dialing conn")
	ErrorAcceptConn = errors.New("Error accepting conn")
	ErrorConnClosed = errors.New("Conn closed")
)