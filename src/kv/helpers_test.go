// vaultreader
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/kv/helpers_test.go

package kv

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"vaultreader/types"
)

// resetGlobals clears the package-level auth/address globals before and
// after each test, since setGlobals mutates them as a side effect and every
// test in this file shares the same package-level state.
func resetGlobals(t *testing.T) {
	t.Helper()
	types.VaultAuthToken = ""
	types.VaultServerAddress = ""
	t.Cleanup(func() {
		types.VaultAuthToken = ""
		types.VaultServerAddress = ""
	})
}

func TestSetGlobalsTokenFlagTakesPrecedence(t *testing.T) {
	resetGlobals(t)
	t.Setenv("VAULT_TOKEN", "token-from-env")
	t.Setenv("VAULT_ADDR", "https://from-env:8200")
	types.VaultAuthToken = "token-from-flag"

	if err := setGlobals(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if types.VaultAuthToken != "token-from-flag" {
		t.Errorf("got %q, want the flag value to win over $VAULT_TOKEN", types.VaultAuthToken)
	}
}

func TestSetGlobalsTokenFallsBackToEnv(t *testing.T) {
	resetGlobals(t)
	t.Setenv("VAULT_TOKEN", "token-from-env")
	t.Setenv("VAULT_ADDR", "https://from-env:8200")

	if err := setGlobals(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if types.VaultAuthToken != "token-from-env" {
		t.Errorf("got %q, want $VAULT_TOKEN value", types.VaultAuthToken)
	}
}

// TestSetGlobalsTokenFallsBackToFile covers the third rung of the resolution
// chain: ~/.vault-token, whitespace-trimmed, used only when neither the
// flag nor $VAULT_TOKEN are set.
func TestSetGlobalsTokenFallsBackToFile(t *testing.T) {
	resetGlobals(t)
	t.Setenv("VAULT_TOKEN", "")
	t.Setenv("VAULT_ADDR", "https://from-env:8200")
	home := t.TempDir()
	t.Setenv("HOME", home)
	if err := os.WriteFile(filepath.Join(home, ".vault-token"), []byte("token-from-file\n"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := setGlobals(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if types.VaultAuthToken != "token-from-file" {
		t.Errorf("got %q, want trimmed contents of ~/.vault-token", types.VaultAuthToken)
	}
}

func TestSetGlobalsTokenMissing(t *testing.T) {
	resetGlobals(t)
	t.Setenv("VAULT_TOKEN", "")
	t.Setenv("HOME", t.TempDir()) // empty dir: no .vault-token file

	err := setGlobals()
	if err == nil {
		t.Fatal("expected an error when no token source is available, got nil")
	}
	if err.Code != types.ErrVaultAuthTokenMissing {
		t.Errorf("error code = %d, want %d (ErrVaultAuthTokenMissing)", err.Code, types.ErrVaultAuthTokenMissing)
	}
}

func TestSetGlobalsAddressFlagTakesPrecedence(t *testing.T) {
	resetGlobals(t)
	t.Setenv("VAULT_TOKEN", "token-from-env")
	t.Setenv("VAULT_ADDR", "https://from-env:8200")
	types.VaultServerAddress = "https://from-flag:8200"

	if err := setGlobals(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if types.VaultServerAddress != "https://from-flag:8200" {
		t.Errorf("got %q, want the flag value to win over $VAULT_ADDR", types.VaultServerAddress)
	}
}

func TestSetGlobalsAddressFallsBackToEnv(t *testing.T) {
	resetGlobals(t)
	t.Setenv("VAULT_TOKEN", "token-from-env")
	t.Setenv("VAULT_ADDR", "https://from-env:8200")

	if err := setGlobals(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if types.VaultServerAddress != "https://from-env:8200" {
		t.Errorf("got %q, want $VAULT_ADDR value", types.VaultServerAddress)
	}
}

func TestSetGlobalsAddressMissing(t *testing.T) {
	resetGlobals(t)
	t.Setenv("VAULT_TOKEN", "token-from-env")
	t.Setenv("VAULT_ADDR", "")

	err := setGlobals()
	if err == nil {
		t.Fatal("expected an error when neither the flag nor $VAULT_ADDR is set, got nil")
	}
	if err.Code != types.ErrVaultServerAddressMissing {
		t.Errorf("error code = %d, want %d (ErrVaultServerAddressMissing)", err.Code, types.ErrVaultServerAddressMissing)
	}
}

// TestClassifyReadError pins down the error-message-substring -> error-code
// mapping that translates vaultLib's wrapped errors into vaultreader's own
// exit codes. Order matters in the source (sealed must be checked before
// the more general "unavailable"), so this also guards against that
// ordering regressing.
func TestClassifyReadError(t *testing.T) {
	cases := []struct {
		name    string
		errText string
		want    int
	}{
		{"sealed", "vault is sealed or unavailable", types.ErrVaultSealed},
		{"unavailable", "Vault is unavailable", types.ErrVaultUnavailable},
		{"connection refused", "dial tcp: connection refused", types.ErrVaultUnavailable},
		{"no such host", "dial tcp: no such host", types.ErrVaultUnavailable},
		{"unauthorized", "unauthorized", types.ErrVaultInvalidAuth},
		{"invalid token", "invalid Vault token", types.ErrVaultInvalidAuth},
		{"permission denied", "permission denied", types.ErrVaultInvalidAuth},
		{"missing path", "secret does not exist", types.ErrInvalidPath},
		{"unrecognized", "something else entirely broke", types.ErrReadSecret},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := classifyReadError(errors.New(tc.errText))
			if got.Code != tc.want {
				t.Errorf("classifyReadError(%q).Code = %d, want %d", tc.errText, got.Code, tc.want)
			}
		})
	}
}

// TestWriteSecretFile checks that writeSecretFile creates the file with
// owner-only permissions and the exact content it was given.
func TestWriteSecretFile(t *testing.T) {
	old := types.SecretOutputFile
	t.Cleanup(func() { types.SecretOutputFile = old })

	dir := t.TempDir()
	types.SecretOutputFile = filepath.Join(dir, "secret.txt")
	types.Quiet = true
	t.Cleanup(func() { types.Quiet = false })

	if err := writeSecretFile([]byte("s3cr3t\n")); err != nil {
		t.Fatalf("writeSecretFile failed: %v", err)
	}

	info, statErr := os.Stat(types.SecretOutputFile)
	if statErr != nil {
		t.Fatal(statErr)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("file permissions = %o, want 600", perm)
	}

	got, readErr := os.ReadFile(types.SecretOutputFile)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(got) != "s3cr3t\n" {
		t.Errorf("file content = %q, want %q", got, "s3cr3t\n")
	}
}

// TestWriteSecretFileError checks that an unwritable destination surfaces as
// ErrWriteFile rather than panicking or being silently swallowed.
func TestWriteSecretFileError(t *testing.T) {
	old := types.SecretOutputFile
	t.Cleanup(func() { types.SecretOutputFile = old })

	// The parent directory doesn't exist, so the write must fail.
	types.SecretOutputFile = filepath.Join(t.TempDir(), "no-such-subdir", "secret.txt")

	err := writeSecretFile([]byte("s3cr3t"))
	if err == nil {
		t.Fatal("expected an error writing to a nonexistent directory, got nil")
	}
	if err.Code != types.ErrWriteFile {
		t.Errorf("error code = %d, want %d (ErrWriteFile)", err.Code, types.ErrWriteFile)
	}
}
