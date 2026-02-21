package checks

import (
	"net"
	"testing"
	"time"
)

func TestCheckTCPSuccess(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		conn, acceptErr := ln.Accept()
		if acceptErr == nil {
			_ = conn.Close()
		}
	}()

	result := CheckTCP(ln.Addr().String(), time.Second)
	if !result.Success {
		t.Fatalf("expected success, got failure: %s", result.Message)
	}
	<-done
}

func TestCheckTCPFailure(t *testing.T) {
	result := CheckTCP("127.0.0.1:1", 100*time.Millisecond)
	if result.Success {
		t.Fatalf("expected failure for closed port")
	}
}
