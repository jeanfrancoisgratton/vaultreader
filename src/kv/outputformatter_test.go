// vaultreader
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/kv/outputformatter_test.go

package kv

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// captureOutputData runs outputData with stdout redirected and returns what
// it printed.
func captureOutputData(t *testing.T, data map[string]interface{}, suppress bool) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	outputData(data, suppress)

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

// TestOutputDataSuppressed checks that suppress=true prints nothing.
func TestOutputDataSuppressed(t *testing.T) {
	if out := captureOutputData(t, map[string]interface{}{"a": 1}, true); out != "" {
		t.Errorf("suppressed call printed %q, want nothing", out)
	}
}

// TestOutputDataWholeSecret checks that every key/value pair is printed as a
// "key: value" line.
func TestOutputDataWholeSecret(t *testing.T) {
	out := captureOutputData(t, map[string]interface{}{"a": 1, "b": "two"}, false)
	for _, want := range []string{"a: 1\n", "b: two\n"} {
		if !strings.Contains(out, want) {
			t.Errorf("output %q does not contain %q", out, want)
		}
	}
}
