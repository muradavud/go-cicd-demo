package checks

import (
	"fmt"
	"net/http"
	"time"
)

type CheckResult struct {
	Success bool
	Message string
}

func CheckHTTP(url string, timeout time.Duration) CheckResult {
	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(url)
	if err != nil {
		return CheckResult{
			Success: false,
			Message: err.Error(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusBadRequest {
		return CheckResult{
			Success: true,
			Message: fmt.Sprintf("status code %d", resp.StatusCode),
		}
	}

	return CheckResult{
		Success: false,
		Message: fmt.Sprintf("status code %d", resp.StatusCode),
	}
}
