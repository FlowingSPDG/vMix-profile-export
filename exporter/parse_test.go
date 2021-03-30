package main

import (
	_ "embed"
	"testing"
)

var (
	//go:embed base.vmix
	baseProfile []byte
)

func TestParseProfile(t *testing.T) {
	p, err := parseProfile(baseProfile)
	if err != nil {
		t.Fatalf("Failed to parse Profile:%v", err)
	}
	t.Logf("Parsed: %#v\n", p)
}
