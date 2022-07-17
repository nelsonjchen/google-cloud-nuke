package cmd

import (
	"testing"
)

func TestIsTrue(t *testing.T) {
	falseStrings := []string{"", "false", "treu", "foo"}
	for _, fs := range falseStrings {
		if IsTrue(fs) {
			t.Fatalf("IsTrue falsely returned 'true' for: %s", fs)
		}
	}

	trueStrings := []string{"true", " true", "true ", " TrUe "}
	for _, ts := range trueStrings {
		if !IsTrue(ts) {
			t.Fatalf("IsTrue falsely returned 'false' for: %s", ts)
		}
	}
}
