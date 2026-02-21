package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunVersion(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{"--version"}, &out, &errOut)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if strings.TrimSpace(out.String()) == "" {
		t.Fatalf("expected non-empty version output")
	}
}

func TestRunNoChecksConfigured(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	code := run([]string{}, &out, &errOut)
	if code != 2 {
		t.Fatalf("expected usage exit code 2, got %d", code)
	}
}

func TestRunRequiredEnvPass(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")

	var out bytes.Buffer
	var errOut bytes.Buffer
	code := run([]string{"--required-env", "DB_HOST"}, &out, &errOut)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(out.String(), "result=true") {
		t.Fatalf("expected success output, got: %s", out.String())
	}
}

func TestRunRequiredEnvFailJSON(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	code := run([]string{"--required-env", "MISSING_VALUE", "--json"}, &out, &errOut)

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(out.String(), "\"success\": false") {
		t.Fatalf("expected JSON failure output, got: %s", out.String())
	}
}
