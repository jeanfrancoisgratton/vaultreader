// vaultreader
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp : 2025/05/31 22:43
// Original filename : src/kv/read.go

package kv

import (
	"encoding/json"
	"fmt"

	"vaultreader/types"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hfjson "github.com/jeanfrancoisgratton/helperFunctions/v5/prettyjson"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	vlr "github.com/jeanfrancoisgratton/vaultLib/kv"
)

// ReadSecrets : Reads a secret fron the Vault secret path
func ReadSecrets(kvengine, path string) *ce.CustomError {
	// Check for required globals
	if err := setGlobals(); err != nil {
		return err
	}

	cfg := vlr.Config{Address: types.VaultServerAddress, Token: types.VaultAuthToken, MountPath: kvengine}
	client, cvlrErr := vlr.NewClient(cfg)
	if cvlrErr != nil {
		return &ce.CustomError{Title: "Error creating vault client", Message: cvlrErr.Error(), Code: types.ErrVaultInit}
	}

	// -f is empty, this means we grab the whole secret
	if types.KVSecretField == "" {
		return allSecrets(client, kvengine, path)
	}

	return singleFieldFromSecret(client, path)
}

// We fetch all the fields of a given secret, optionally rendering it in JSON
func allSecrets(c *vlr.Client, kvengine, path string) *ce.CustomError {
	var secret *vlr.Secret
	var sErr error

	if secret, sErr = c.ReadSecret(path,
		vlr.ReadOptions{Version: types.KVSecretVersion, FallbackToLatestAvailable: true}); sErr != nil {
		return classifyReadError(sErr)
	}

	if types.OutputFormat == "json" {
		payload, err := json.MarshalIndent(secret.Data, "", "  ")
		if err != nil {
			return &ce.CustomError{Title: "Error serializing secret", Message: err.Error(), Code: types.ErrExtractData}
		}
		if e := hfjson.Print(payload); e != nil {
			return &ce.CustomError{Title: "Unable to render secret's payload", Message: e.Error(), Code: types.ErrExtractData}
		}
		return nil
	}
	return outputData(secret.Data, types.Quiet)
}

func singleFieldFromSecret(c *vlr.Client, path string) *ce.CustomError {
	value, err := c.ReadSecretField(path, types.KVSecretField, types.KVSecretVersion)
	if err != nil {
		return classifyReadError(err)
	}

	if types.Quiet {
		fmt.Printf("%v\n", value)
	} else {
		fmt.Println(types.KVSecretField + " : " + hftx.Green(fmt.Sprintf("%v", value)))
	}
	return nil
}
