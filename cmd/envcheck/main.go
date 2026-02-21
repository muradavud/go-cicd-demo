package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/muradavid/go-cicd-demo/internal/checks"
	"github.com/muradavid/go-cicd-demo/internal/version"
)

type stringListFlag []string

func (s *stringListFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *stringListFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

type checkResult struct {
	Name    string `json:"name"`
	Target  string `json:"target"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type runResult struct {
	Success  bool          `json:"success"`
	Checks   []checkResult `json:"checks"`
	Duration string        `json:"duration"`
	Version  string        `json:"version"`
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("envcheck", flag.ContinueOnError)
	fs.SetOutput(stderr)

	requiredEnvFlag := fs.String("required-env", "", "comma-separated list of required environment variables")
	timeout := fs.Duration("timeout", 5*time.Second, "timeout for network checks, e.g. 3s")
	jsonOutput := fs.Bool("json", false, "emit machine-readable JSON output")
	showVersion := fs.Bool("version", false, "print build version")

	var httpChecks stringListFlag
	var tcpChecks stringListFlag
	fs.Var(&httpChecks, "check-http", "HTTP endpoint to check; repeatable")
	fs.Var(&tcpChecks, "check-tcp", "TCP address to check (host:port); repeatable")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *showVersion {
		fmt.Fprintln(stdout, version.Full())
		return 0
	}

	requiredEnv := splitCSV(*requiredEnvFlag)
	httpTargets := splitAndFlatten(httpChecks)
	tcpTargets := splitAndFlatten(tcpChecks)
	if len(requiredEnv) == 0 && len(httpTargets) == 0 && len(tcpTargets) == 0 {
		fmt.Fprintln(stderr, "at least one check must be configured")
		fs.Usage()
		return 2
	}

	start := time.Now()
	results := make([]checkResult, 0, len(requiredEnv)+len(httpTargets)+len(tcpTargets))
	failures := 0

	envResults := checks.ValidateRequiredEnv(requiredEnv)
	for _, r := range envResults {
		results = append(results, checkResult{
			Name:    "env",
			Target:  r.Key,
			Success: r.Success,
			Message: r.Message,
		})
		if !r.Success {
			failures++
		}
	}

	for _, target := range httpTargets {
		r := checks.CheckHTTP(target, *timeout)
		results = append(results, checkResult{
			Name:    "http",
			Target:  target,
			Success: r.Success,
			Message: r.Message,
		})
		if !r.Success {
			failures++
		}
	}

	for _, target := range tcpTargets {
		r := checks.CheckTCP(target, *timeout)
		results = append(results, checkResult{
			Name:    "tcp",
			Target:  target,
			Success: r.Success,
			Message: r.Message,
		})
		if !r.Success {
			failures++
		}
	}

	output := runResult{
		Success:  failures == 0,
		Checks:   results,
		Duration: time.Since(start).String(),
		Version:  version.Full(),
	}

	if *jsonOutput {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(output); err != nil {
			fmt.Fprintf(stderr, "failed to encode JSON output: %v\n", err)
			return 2
		}
	} else {
		for _, r := range output.Checks {
			status := "PASS"
			if !r.Success {
				status = "FAIL"
			}
			fmt.Fprintf(stdout, "[%s] check=%s target=%s msg=%q\n", status, r.Name, r.Target, r.Message)
		}
		fmt.Fprintf(stdout, "result=%t checks=%d duration=%s version=%s\n", output.Success, len(output.Checks), output.Duration, output.Version)
	}

	if output.Success {
		return 0
	}
	return 1
}

func splitAndFlatten(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		out = append(out, splitCSV(value)...)
	}
	return out
}

func splitCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		clean := strings.TrimSpace(p)
		if clean != "" {
			out = append(out, clean)
		}
	}
	return out
}
