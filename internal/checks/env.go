package checks

import "os"

type EnvResult struct {
	Key     string
	Success bool
	Message string
}

func ValidateRequiredEnv(required []string) []EnvResult {
	results := make([]EnvResult, 0, len(required))

	for _, key := range required {
		value, exists := os.LookupEnv(key)
		if !exists || value == "" {
			results = append(results, EnvResult{
				Key:     key,
				Success: false,
				Message: "environment variable is missing or empty",
			})
			continue
		}

		results = append(results, EnvResult{
			Key:     key,
			Success: true,
			Message: "present",
		})
	}

	return results
}
