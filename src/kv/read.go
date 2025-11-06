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

	hfl "github.com/jeanfrancoisgratton/helperFunctions/v3/logging"
)

// ReadSecrets : Reads a secret fron the Vault secret path
func ReadSecrets(path string) int {
	hfl.Debugf("Starting ReadSecrets for path=%s", path)

	// Check for required globals
	if err := setGlobals(); err != 0 {
		return err
	}

	// Create Vault client
	cfg := &api.Config{Address: types.VaultServerAddress}
	client, err := api.NewClient(cfg)
	if err != nil {
		hfl.Errorf("Vault client creation failed: %v", err)
		return types.ErrVaultInit
	}
	client.SetToken(types.VaultAuthToken)

	dataPath := fmt.Sprintf("%s/data/%s", strings.Trim(types.KVEngineMountPath, "/"), strings.Trim(path, "/"))
	metaPath := fmt.Sprintf("%s/metadata/%s", strings.Trim(types.KVEngineMountPath, "/"), strings.Trim(path, "/"))

	// Pre-check if metadata exists
	_, metaErr := client.Logical().Read(metaPath)
	if metaErr != nil {
		if strings.Contains(metaErr.Error(), "connection refused") || strings.Contains(metaErr.Error(), "no such host") {
			hfl.Errorf("Vault service unavailable: %v", metaErr)
			return types.ErrVaultUnavailable
		}
		if strings.Contains(metaErr.Error(), "permission denied") || strings.Contains(metaErr.Error(), "unauthorized") {
			hfl.Errorf("Invalid Vault token or unauthorized: %v", metaErr)
			return types.ErrVaultInvalidAuth
		}
		if strings.Contains(metaErr.Error(), "server is sealed") {
			hfl.Errorf("Vault is sealed: %v", metaErr)
			return types.ErrVaultSealed
		}

		hfl.Errorf("Secret path does not exist or metadata read failed: %s", path)
		return types.ErrInvalidPath
	}

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
				hfl.Errorf("Fallback version detection failed: %v", ferr)
				return types.ErrReadSecret
			}
			if ver == 0 {
				hfl.Errorf("All versions destroyed")
				return types.ErrReadSecret
			}
			hfl.Infof("Using fallback available version %d", ver)
			secret, err = client.Logical().ReadWithData(dataPath, map[string][]string{
				"version": {fmt.Sprintf("%d", ver)},
			})
		}
	}

	if err != nil {
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "no such host") {
			hfl.Errorf("Vault service unavailable: %v", err)
			return types.ErrVaultUnavailable
		}
		if strings.Contains(err.Error(), "permission denied") || strings.Contains(err.Error(), "unauthorized") {
			hfl.Errorf("Invalid Vault token or unauthorized: %v", err)
			return types.ErrVaultInvalidAuth
		}
		if strings.Contains(err.Error(), "server is sealed") {
			hfl.Errorf("Vault is sealed: %v", err)
			return types.ErrVaultSealed
		}

		hfl.Errorf("Secret read failed: %v", err)
		return types.ErrReadSecret
	}

	if secret == nil {
		hfl.Errorf("Secret read returned nil")
		return types.ErrReadSecret
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		hfl.Errorf("Secret format invalid")
		return types.ErrExtractData
	}

	return outputData(data, types.Quiet)
}
