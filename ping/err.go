package ping

import "errors"

var (
	TimeoutErr         = errors.New("timeout")
	HostUnreachableErr = errors.New("host is unreachable")
)
