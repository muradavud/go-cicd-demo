package checks

import "testing"

func TestValidateRequiredEnv(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("EMPTY_ALLOWED", "")

	results := ValidateRequiredEnv([]string{"DB_HOST", "MISSING", "EMPTY_ALLOWED"})
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	if !results[0].Success {
		t.Fatalf("expected DB_HOST to pass")
	}
	if results[1].Success {
		t.Fatalf("expected MISSING to fail")
	}
	if results[2].Success {
		t.Fatalf("expected EMPTY_ALLOWED to fail")
	}
}
