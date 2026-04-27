package errs

import "testing"

func TestToHTTP(t *testing.T) {
	status, code, message := ToHTTP(New(CodeUnauthorized, "unauthorized"))
	if status != 401 || code != int(CodeUnauthorized) || message != "unauthorized" {
		t.Fatalf("unexpected mapping: %d %d %s", status, code, message)
	}
}
