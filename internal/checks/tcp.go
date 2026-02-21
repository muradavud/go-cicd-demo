package checks

import (
	"net"
	"time"
)

func CheckTCP(address string, timeout time.Duration) CheckResult {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return CheckResult{
			Success: false,
			Message: err.Error(),
		}
	}
	_ = conn.Close()

	return CheckResult{
		Success: true,
		Message: "connection successful",
	}
}
