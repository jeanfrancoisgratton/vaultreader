// vaultreader
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/kv/outputformatter_test.go

package kv

import (
	"testing"

	"vaultreader/types"
)

func resetOutputGlobals(t *testing.T) {
	t.Helper()
	oldField := types.KVSecretField
	oldFormat := types.OutputFormat
	t.Cleanup(func() {
		types.KVSecretField = oldField
		types.OutputFormat = oldFormat
	})
}

// TestOutputDataFieldNotFound checks the one code path in outputData that
// returns an error regardless of the suppress flag: a requested field that
// isn't present in the secret's data.
func TestOutputDataFieldNotFound(t *testing.T) {
	resetOutputGlobals(t)
	types.KVSecretField = "missing-field"

	err := outputData(map[string]interface{}{"present-field": "value"}, true)
	if err == nil {
		t.Fatal("expected an error for a missing field, got nil")
	}
	if err.Code != types.ErrFieldNotFound {
		t.Errorf("error code = %d, want %d (ErrFieldNotFound)", err.Code, types.ErrFieldNotFound)
	}
}

// TestOutputDataSuppressed exercises the suppress=true path for both the
// whole-secret and single-field cases: neither should return an error.
func TestOutputDataSuppressed(t *testing.T) {
	resetOutputGlobals(t)

	types.KVSecretField = ""
	if err := outputData(map[string]interface{}{"a": 1}, true); err != nil {
		t.Errorf("whole-secret suppressed call returned an error: %v", err)
	}

	types.KVSecretField = "a"
	if err := outputData(map[string]interface{}{"a": 1}, true); err != nil {
		t.Errorf("single-field suppressed call returned an error: %v", err)
	}
}

// TestOutputDataFieldFound checks the non-suppressed, field-present path
// returns no error for both text and JSON output formats.
func TestOutputDataFieldFound(t *testing.T) {
	resetOutputGlobals(t)
	types.KVSecretField = "a"

	types.OutputFormat = "text"
	if err := outputData(map[string]interface{}{"a": 1}, false); err != nil {
		t.Errorf("text output returned an error: %v", err)
	}

	types.OutputFormat = "json"
	if err := outputData(map[string]interface{}{"a": 1}, false); err != nil {
		t.Errorf("json output returned an error: %v", err)
	}
}

// TestOutputDataWholeSecret checks the non-suppressed, whole-secret path
// (no field selected) returns no error for both text and JSON.
func TestOutputDataWholeSecret(t *testing.T) {
	resetOutputGlobals(t)
	types.KVSecretField = ""

	types.OutputFormat = "text"
	if err := outputData(map[string]interface{}{"a": 1, "b": 2}, false); err != nil {
		t.Errorf("text output returned an error: %v", err)
	}

	types.OutputFormat = "json"
	if err := outputData(map[string]interface{}{"a": 1, "b": 2}, false); err != nil {
		t.Errorf("json output returned an error: %v", err)
	}
}
