package customError

import "errors"

var (
	ErrorDialConn   = errors.New("Error when dialing conn")
	ErrorAcceptConn = errors.New("Error accepting conn")
	ErrorConnClosed = errors.New("Conn closed")
)

var (
	ERR_GET_CONN_TIMEOUT      = errors.New("Error get conn timeout")
	ERR_FREE_CONN_CHAN_CLOSED = errors.New("Free conn channel error")
	ERR_CONN_POOL_QUEUE_FULL  = errors.New("Error queue is full")
)
