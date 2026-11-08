package main

import "testing"

func TestVersionString(t *testing.T) {
	old := version
	t.Cleanup(func() { version = old })

	version = "v0.1.0"
	if got := versionString(); got != "v0.1.0" {
		t.Fatalf("got %q", got)
	}
}
