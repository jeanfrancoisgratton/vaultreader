// vaultreader
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp : 2025/06/04 10:04
// Original filename : src/kv/helpers.go

package kv

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/api"

	ce "github.com/jeanfrancoisgratton/customError/v3"
)

// findLatestAvailableVersion :
// In a KV engine version 2+, multiple versions of a secret may exist.
// We may need to find the latest recorded version
func findLatestAvailableVersion(client *api.Client, metaPath string) (int, *ce.CustomError) {
	meta, err := client.Logical().Read(metaPath)
	if err != nil || meta == nil {
		// meta can be nil with a nil err (path simply has no metadata), so we
		// must not blindly dereference err here.
		message := "No metadata found at " + metaPath
		if err != nil {
			message = err.Error()
		}
		return 0, &ce.CustomError{Title: "Unable to fetch metadata", Message: message, Code: ErrReadSecret}
	}

	rawVersions, ok := meta.Data["versions"].(map[string]interface{})
	if !ok {
		// err is guaranteed nil at this point; the failure is the missing/
		// malformed "versions" key, not a read error.
		return 0, &ce.CustomError{Title: "Version metadata not found",
			Message: "Metadata at " + metaPath + " has no usable \"versions\" map", Code: ErrExtractData}
	}

	var available []int
	for verStr, vmetaAny := range rawVersions {
		vmeta, ok := vmetaAny.(map[string]interface{})
		if !ok {
			continue
		}
		if destroyed, _ := vmeta["destroyed"].(bool); destroyed {
			continue
		}
		if ver, err := strconv.Atoi(verStr); err == nil {
			available = append(available, ver)
		}
	}

	if len(available) == 0 {
		return 0, nil
	}

	sort.Sort(sort.Reverse(sort.IntSlice(available)))
	return available[0], nil
}

// setGlobals :
// Setting the VAULT_TOKEN and VAULT_ADDR values, either from the env vars or command-line flags
func setGlobals() *ce.CustomError {
	// We check if we have a valid token value, if not we exit
	if VaultAuthToken == "" {
		VaultAuthToken = os.Getenv("VAULT_TOKEN")
	}
	// the -t flag and VAULT_TOKEN env var are not set, we check if there is a $HOME/.vault-token file
	if VaultAuthToken == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			data, err := os.ReadFile(filepath.Join(homeDir, ".vault-token"))
			if err == nil {
				VaultAuthToken = strings.TrimSpace(string(data))
			}
		}
	}
	// -t and VAULT_TOKEN are not set, $HOME/.vault-token is missing or empty
	if VaultAuthToken == "" {
		errInfo := ErrorMessages[ErrVaultAuthTokenMissing]
		title := fmt.Sprintf("[%s] Vault token is missing", errInfo.Int2StringCode)
		message := fmt.Sprintf("Neither the $VAULT_TOKEN variable, the -t flag or the ~%s/.vault-token file were set.",
			filepath.Base(os.Getenv("HOME")))
		return &ce.CustomError{Title: title, Message: message, Code: ErrVaultAuthTokenMissing}
	}
	// ok, so we have a token, let's now check if we have a valid vault server address, be it
	// in an environment variable or with the -a flag
	if VaultServerAddress == "" {
		VaultServerAddress = os.Getenv("VAULT_ADDR")
	}
	if VaultServerAddress == "" {
		title := "Vault address is missing"
		message := "Neither the $VAULT_ADDR variable or the -a flag were set"
		return &ce.CustomError{Title: title, Message: message, Code: ErrVaultServerAddressMissing}
	}
	return nil
}

// writeSecretFile writes a secret's rendered content to SecretOutputFile
// with owner-only permissions, since the file may contain sensitive material.
func writeSecretFile(content []byte) *ce.CustomError {
	if err := os.WriteFile(SecretOutputFile, content, 0600); err != nil {
		return &ce.CustomError{Title: ErrorMessages[ErrWriteFile].Msg, Message: err.Error(), Code: ErrWriteFile}
	}
	if !Quiet {
		fmt.Printf("Secret written to %s\n", SecretOutputFile)
	}
	return nil
}

// classifyReadError maps an error returned by vaultLib's KV read path onto a
// CustomError carrying the matching vaultreader error code, so that Die() can
// translate it into a meaningful, POSIX-safe exit status (e.g. a sealed vault
// exits with ErrVaultSealed rather than an indistinguishable generic 1).
//
// vaultLib normalises its failures into recognisable wrapped strings; we match
// on those. Order matters: the sealed case must be tested before the more
// general "unavailable" one, since vaultLib reports "sealed or unavailable".
func classifyReadError(err error) *ce.CustomError {
	msg := err.Error()
	var code int
	switch {
	case strings.Contains(msg, "sealed"):
		code = ErrVaultSealed
	case strings.Contains(msg, "unavailable"), strings.Contains(msg, "connection refused"),
		strings.Contains(msg, "no such host"):
		code = ErrVaultUnavailable
	case strings.Contains(msg, "unauthorized"), strings.Contains(msg, "invalid Vault token"),
		strings.Contains(msg, "permission denied"):
		code = ErrVaultInvalidAuth
	case strings.Contains(msg, "does not exist"),
		strings.Contains(msg, "read secret failed: secret ") && strings.Contains(msg, "was not found"):
		code = ErrInvalidPath
	// Matched on vaultLib's full prefix: a bare "not found" would also catch
	// "secret ... was not found", which is a path problem, not a field one.
	case strings.Contains(msg, "read secret field failed: field "):
		code = ErrFieldNotFound
	default:
		code = ErrReadSecret
	}
	return &ce.CustomError{Title: ErrorMessages[code].Msg, Message: msg, Code: code}
}
