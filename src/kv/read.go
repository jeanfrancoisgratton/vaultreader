// vaultreader
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp : 2025/05/31 22:43
// Original filename : src/kv/read.go

package kv

import (
	"fmt"
	"strings"

	hfl "github.com/jeanfrancoisgratton/helperFunctions/v5/logging"
	"vaultreader/types"

	"github.com/hashicorp/vault/api"
	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
)

// ReadSecrets : Reads a secret fron the Vault secret path
// FIXME: split this function is smaller functions
func ReadSecrets(path string) *ce.CustomError {
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
			fmt.Println(hftx.SkullBonesSign(title + ": " + message))
		}
		return &ce.CustomError{Title: title, Message: message, Code: types.ErrVaultAuthTokenMissing}
	}
	client.SetToken(types.VaultAuthToken)

	dataPath := fmt.Sprintf("%s/data/%s", strings.Trim(types.KVEngineMountPath, "/"), strings.Trim(path, "/"))
	metaPath := fmt.Sprintf("%s/metadata/%s", strings.Trim(types.KVEngineMountPath, "/"), strings.Trim(path, "/"))

	// Pre-check metadata
	_, metaErr := client.Logical().Read(metaPath)

	if metaErr != nil {
		// vault unavail
		if strings.Contains(metaErr.Error(), "connection refused") || strings.Contains(metaErr.Error(), "no such host") {
			title := "Vault service unavailable"
			message := metaErr.Error()
			if !types.Quiet {
				fmt.Println(hftx.SkullBonesSign(title + ": " + message))
			}
			return &ce.CustomError{Title: title, Message: message, Code: types.ErrVaultUnavailable}
		}
		// 403 (not authorized)
		if strings.Contains(metaErr.Error(), "permission denied") || strings.Contains(metaErr.Error(), "unauthorized") {
			title := "Invalid Vault token or unauthorized"
			message := metaErr.Error()
			if !types.Quiet {
				fmt.Println(hftx.SkullBonesSign(title + ": " + message))
			}
			return &ce.CustomError{Title: title, Message: message, Code: types.ErrVaultInvalidAuth}
		}
		// vault is sealed
		if strings.Contains(metaErr.Error(), "Vault is sealed") {
			title := "Vault is sealed"
			message := metaErr.Error()
			if !types.Quiet {
				fmt.Println(hftx.SkullBonesSign(title + ": " + message))
			}
			return &ce.CustomError{Title: title, Message: message, Code: types.ErrVaultSealed}
		}

		// If we're here, this means that either the secret path or metadata do not exist
		title := "Secret path does not exist or metadata read failed"
		message := metaErr.Error()
		if !types.Quiet {
			fmt.Println(hftx.SkullBonesSign(title + ": " + message))
		}

		return &ce.CustomError{Title: title, Message: message, Code: types.ErrVaultSealed}
	}

	// so, ... we good ?
	var secret *api.Secret

	if types.KVSecretVersion > 0 {
		secret, err = client.Logical().ReadWithData(dataPath, map[string][]string{
			"version": {fmt.Sprintf("%d", types.KVSecretVersion)},
		})
	} else {
		secret, err = client.Logical().Read(dataPath)
		if err == nil && secret == nil {
			ver, ferr := findLatestAvailableVersion(client, metaPath)
			if ferr != nil {
				ferr.Code = types.ErrReadSecret
				if !types.Quiet {
					fmt.Println(hftx.SkullBonesSign(ferr.Error()))
				}
				return ferr
			}
			if ver == 0 {
				if !types.Quiet {
					fmt.Println(hftx.SkullBonesSign(" All the secret's versions were destroyed"))
				}
				return &ce.CustomError{Title: "ReadSecret failed",
					Message: "All the secret's versions were destroyed",
					Code:    types.ErrReadSecret}
			}
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
				fmt.Println(hftx.SkullBonesSign(title + ": " + message))
			}
			err := ce.CustomError{Title: title, Message: message, Code: types.ErrVaultUnavailable}
			hfl.Errorf(err.ErrorNoColor())
			return &err
		}
		if strings.Contains(err.Error(), "permission denied") || strings.Contains(err.Error(), "unauthorized") {
			title := "Invalid Vault token or unauthorized"
			message := err.Error()
			if !types.Quiet {
				fmt.Println(hftx.SkullBonesSign(title + ": " + message))
			}
			cerr := ce.CustomError{Title: title, Message: message, Code: types.ErrVaultInvalidAuth}
			hfl.Errorf(cerr.ErrorNoColor())
			return &cerr
		}
		if strings.Contains(err.Error(), "server is sealed") {
			title := "Vault is sealed"
			message := err.Error()
			if !types.Quiet {
				fmt.Println(hftx.SkullBonesSign(title + ": " + message))
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
			fmt.Println(hftx.SkullBonesSign(title + ": " + message))
		}
		return &cerr
	}

	if secret == nil {
		title := "ReadSecret failed"
		message := "Secret read returned nil"
		code := types.ErrReadSecret
		cerr := ce.CustomError{Title: title, Message: message, Code: code}
		if !types.Quiet {
			fmt.Println(hftx.SkullBonesSign(title + ": " + message))
		}
		return &cerr
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		title := "ReadSecret failed"
		message := "Secret format is invalid"
		code := types.ErrExtractData
		cerr := ce.CustomError{Title: title, Message: message, Code: code}
		if !types.Quiet {
			fmt.Println(hftx.SkullBonesSign(title + ": " + message))
		}
		return &cerr
	}

	return outputData(data, types.Quiet)
}
