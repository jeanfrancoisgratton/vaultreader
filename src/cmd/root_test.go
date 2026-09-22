// vaultreader
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/root_test.go

package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// TestVersionCommandOutput checks that `vaultreader version` prints the
// ldflags-injected buildVersion/buildDate (rather than, e.g., silently
// falling back to the "dev"/"unknown" zero values) — this is the command's
// only behavior, and it needs no Vault server.
func TestVersionCommandOutput(t *testing.T) {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	versionCmd.Run(versionCmd, nil)

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "vaultreader v"+buildVersion) {
		t.Errorf("output %q does not contain the build version %q", out, buildVersion)
	}
	if !strings.Contains(out, buildDate) {
		t.Errorf("output %q does not contain the build date %q", out, buildDate)
	}
}

// TestFileFlagHasNoShorthandCollision guards the fix applied when --file
// was added: it must NOT take -f, since -f is already bound to --field on
// the root command. Regressing this would make one of the two flags
// silently shadow the other.
func TestFileFlagHasNoShorthandCollision(t *testing.T) {
	fieldFlag := rootCmd.PersistentFlags().ShorthandLookup("f")
	if fieldFlag == nil {
		t.Fatal("expected -f to be registered as a shorthand on the root command")
	}
	if fieldFlag.Name != "field" {
		t.Errorf("-f is bound to %q, want it bound to \"field\"", fieldFlag.Name)
	}

	fileFlag := rootCmd.PersistentFlags().Lookup("file")
	if fileFlag == nil {
		t.Fatal("expected --file to be registered on the root command")
	}
	if fileFlag.Shorthand != "" {
		t.Errorf("--file has shorthand -%s, want no shorthand (it would collide with -f/--field)", fileFlag.Shorthand)
	}
}
