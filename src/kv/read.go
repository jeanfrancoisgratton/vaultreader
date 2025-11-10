// vaultreader
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp : 2025/05/31 22:43
// Original filename : src/kv/read.go

package kv

import (
	"fmt"
	"strings"
	"vaultreader/types"

	"github.com/hashicorp/vault/api"
	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v3/terminalfx"

	hfl "github.com/jeanfrancoisgratton/helperFunctions/v3/logging"
)

// ReadSecrets : Reads a secret fron the Vault secret path
// FIXME: split this function is smaller functions
func ReadSecrets(path string) *ce.CustomError {
	hfl.Debugf("Starting ReadSecrets for path=%s", path)

	// Check for required globals
	if err := setGlobals(); err != nil {
		return err
	}

	// Create Vault client
	cfg := &api.Config{Address: types.VaultServerAddress}
	client, err := api.NewClient(cfg)
	if err != nil {
		title := "Vault client creation failed"
		message := err.Error()
		if !types.Quiet {
			fmt.Println(hftx.FatalSkullBonesGlyph(fmt.Sprintf("%s:\n%s", title, message)))
		}
		cerr := ce.CustomError{Title: title, Message: message, Code: types.ErrVaultAuthTokenMissing}
		hfl.Errorf(cerr.ErrorNoColor())
		return &cerr
	}
	client.SetToken(types.VaultAuthToken)

	dataPath := fmt.Sprintf("%s/data/%s", strings.Trim(types.KVEngineMountPath, "/"), strings.Trim(path, "/"))
	metaPath := fmt.Sprintf("%s/metadata/%s", strings.Trim(types.KVEngineMountPath, "/"), strings.Trim(path, "/"))
	hfl.Debugf("Data path is ", dataPath)
	hfl.Debugf("Metadata path is ", metaPath)

	// Pre-check metadata
	_, metaErr := client.Logical().Read(metaPath)

	if metaErr != nil {
		// vault unavail
		if strings.Contains(metaErr.Error(), "connection refused") || strings.Contains(metaErr.Error(), "no such host") {
			title := "Vault service unavailable"
			message := metaErr.Error()
			if !types.Quiet {
				fmt.Println(hftx.FatalSkullBonesGlyph(fmt.Sprintf("%s:\n%s", title, message)))
			}
			cerr := ce.CustomError{Title: title, Message: message, Code: types.ErrVaultUnavailable}
			hfl.Errorf(cerr.ErrorNoColor())
			return &cerr
		}
		// 403 (not authorized)
		if strings.Contains(metaErr.Error(), "permission denied") || strings.Contains(metaErr.Error(), "unauthorized") {
			title := "Invalid Vault token or unauthorized"
			message := metaErr.Error()
			if !types.Quiet {
				fmt.Println(hftx.FatalSkullBonesGlyph(fmt.Sprintf("%s:\n%s", title, message)))
			}
			cerr := ce.CustomError{Title: title, Message: message, Code: types.ErrVaultInvalidAuth}
			hfl.Errorf(cerr.ErrorNoColor())
			return &cerr
		}
		// vault is sealed
		if strings.Contains(metaErr.Error(), "server is sealed") {
			title := "Vault is sealed"
			message := metaErr.Error()
			if !types.Quiet {
				fmt.Println(hftx.FatalSkullBonesGlyph(fmt.Sprintf("%s:\n%s", title, message)))
			}
			cerr := ce.CustomError{Title: title, Message: message, Code: types.ErrVaultSealed}
			hfl.Errorf(cerr.ErrorNoColor())
			return &cerr
		}

		// If we're here, this means that either the secret path or metadata do not exist
		title := "Secret path does not exist or metadata read failed"
		message := metaErr.Error()
		if !types.Quiet {
			fmt.Println(hftx.FatalSkullBonesGlyph(fmt.Sprintf("%s:\n%s", title, message)))
		}
		cerr := ce.CustomError{Title: title, Message: message, Code: types.ErrVaultSealed}
		hfl.Errorf(cerr.ErrorNoColor())
		return &cerr
	}

	// so, ... we good ?
	var secret *api.Secret

	if types.KVSecretVersion > 0 {
		hfl.Infof("Fetching version %d from %s", types.KVSecretVersion, dataPath)
		hfl.Debugf("Explicit secret read: mount=%s path=%s version=%d", types.KVEngineMountPath, path, types.KVSecretVersion)
		secret, err = client.Logical().ReadWithData(dataPath, map[string][]string{
			"version": {fmt.Sprintf("%d", types.KVSecretVersion)},
		})
	} else {
		hfl.Debugf("Attempting latest version read: mount=%s path=%s", types.KVEngineMountPath, path)
		secret, err = client.Logical().Read(dataPath)
		if err == nil && secret == nil {
			hfl.Infof("Latest version missing; trying to find fallback version")
			ver, ferr := findLatestAvailableVersion(client, metaPath)
			if ferr != nil {
				ferr.Code = types.ErrReadSecret
				if !types.Quiet {
					fmt.Println(hftx.FatalSkullBonesGlyph(ferr.Error()))
				}
				hfl.Errorf(ferr.ErrorNoColor())
				return ferr
			}
			if ver == 0 {
				if !types.Quiet {
					fmt.Println(hftx.FatalSkullBonesGlyph("All the secret's versions were destroyed"))
				}
				cerr := ce.CustomError{Title: "ReadSecret failed",
					Message: "All the secret's versions were destroyed",
					Code:    types.ErrReadSecret}
				hfl.Errorf(cerr.ErrorNoColor())
				return &cerr
			}
			hfl.Infof("Using fallback available version %d", ver)
			secret, err = client.Logical().ReadWithData(dataPath, map[string][]string{
				"version": {fmt.Sprintf("%d", ver)},
			})
		}
	}

	if err != nil {
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "no such host") {
			title := "Vault service unavailable"
			message := err.Error()
			if !types.Quiet {
				fmt.Println(hftx.FatalSkullBonesGlyph(fmt.Sprintf("%s:\n%s", title, message)))
			}
			err := ce.CustomError{Title: title, Message: message, Code: types.ErrVaultUnavailable}
			hfl.Errorf(err.ErrorNoColor())
			return &err
		}
		if strings.Contains(err.Error(), "permission denied") || strings.Contains(err.Error(), "unauthorized") {
			title := "Invalid Vault token or unauthorized"
			message := err.Error()
			if !types.Quiet {
				fmt.Println(hftx.FatalSkullBonesGlyph(fmt.Sprintf("%s:\n%s", title, message)))
			}
			cerr := ce.CustomError{Title: title, Message: message, Code: types.ErrVaultInvalidAuth}
			hfl.Errorf(cerr.ErrorNoColor())
			return &cerr
		}
		if strings.Contains(err.Error(), "server is sealed") {
			title := "Vault is sealed"
			message := err.Error()
			if !types.Quiet {
				fmt.Println(hftx.FatalSkullBonesGlyph(fmt.Sprintf("%s:\n%s", title, message)))
			}
			cerr := ce.CustomError{Title: title, Message: message, Code: types.ErrVaultSealed}
			hfl.Errorf(cerr.ErrorNoColor())
			return &cerr
		}
		title := "ReadSecret failed"
		message := err.Error()
		code := types.ErrReadSecret
		cerr := ce.CustomError{Title: title, Message: message, Code: code}
		if !types.Quiet {
			fmt.Println(hftx.FatalSkullBonesGlyph(fmt.Sprintf("%s:\n%s", title, message)))
		}
		hfl.Errorf(cerr.ErrorNoColor())
		return &cerr
	}

	if secret == nil {
		title := "ReadSecret failed"
		message := "Secret read returned nil"
		code := types.ErrReadSecret
		cerr := ce.CustomError{Title: title, Message: message, Code: code}
		if !types.Quiet {
			fmt.Println(hftx.FatalSkullBonesGlyph(fmt.Sprintf("%s:\n%s", title, message)))
		}
		hfl.Errorf(cerr.ErrorNoColor())
		return &cerr
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		title := "ReadSecret failed"
		message := "Secret format is invalid"
		code := types.ErrExtractData
		cerr := ce.CustomError{Title: title, Message: message, Code: code}
		if !types.Quiet {
			fmt.Println(hftx.FatalSkullBonesGlyph(fmt.Sprintf("%s:\n%s", title, message)))
		}
		hfl.Errorf(cerr.ErrorNoColor())
		return &cerr
	}

	return outputData(data, types.Quiet)
}
