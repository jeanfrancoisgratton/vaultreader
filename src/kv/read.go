// vaultreader
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp : 2025/05/31 22:43
// Original filename : src/kv/read.go

package kv

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"strings"
	"vaultreader/logging"
	"vaultreader/types"
)

func ReadSecrets(path string) int {
	logging.Debugf("Starting ReadSecrets for path=%s", path)

	// Check for required globals
	if err := setGlobals(); err != 0 {
		return err
	}

	// Create Vault client
	cfg := &api.Config{Address: types.VaultServerAddress}
	client, err := api.NewClient(cfg)
	if err != nil {
		logging.Errorf("Vault client creation failed: %v", err)
		return types.ErrVaultInit
	}
	client.SetToken(types.VaultAuthToken)

	dataPath := fmt.Sprintf("%s/data/%s", strings.Trim(types.KVEngineMountPath, "/"), strings.Trim(path, "/"))
	metaPath := fmt.Sprintf("%s/metadata/%s", strings.Trim(types.KVEngineMountPath, "/"), strings.Trim(path, "/"))

	// Pre-check if metadata exists
	_, metaErr := client.Logical().Read(metaPath)
	if metaErr != nil {
		if strings.Contains(metaErr.Error(), "connection refused") || strings.Contains(metaErr.Error(), "no such host") {
			logging.Errorf("Vault service unavailable: %v", metaErr)
			return types.ErrVaultUnavailable
		}
		if strings.Contains(metaErr.Error(), "permission denied") || strings.Contains(metaErr.Error(), "unauthorized") {
			logging.Errorf("Invalid Vault token or unauthorized: %v", metaErr)
			return types.ErrVaultInvalidAuth
		}
		if strings.Contains(metaErr.Error(), "server is sealed") {
			logging.Errorf("Vault is sealed: %v", metaErr)
			return types.ErrVaultSealed
		}

		logging.Errorf("Secret path does not exist or metadata read failed: %s", path)
		return types.ErrInvalidPath
	}

	var secret *api.Secret

	if types.KVSecretVersion > 0 {
		logging.Infof("Fetching version %d from %s", types.KVSecretVersion, dataPath)
		logging.Debugf("Explicit secret read: mount=%s path=%s version=%d", types.KVEngineMountPath, path, types.KVSecretVersion)
		secret, err = client.Logical().ReadWithData(dataPath, map[string][]string{
			"version": {fmt.Sprintf("%d", types.KVSecretVersion)},
		})
	} else {
		logging.Debugf("Attempting latest version read: mount=%s path=%s", types.KVEngineMountPath, path)
		secret, err = client.Logical().Read(dataPath)
		if err == nil && secret == nil {
			logging.Infof("Latest version missing; trying to find fallback version")
			ver, ferr := findLatestAvailableVersion(client, metaPath)
			if ferr != nil {
				logging.Errorf("Fallback version detection failed: %v", ferr)
				return types.ErrReadSecret
			}
			if ver == 0 {
				logging.Errorf("All versions destroyed")
				return types.ErrReadSecret
			}
			logging.Infof("Using fallback available version %d", ver)
			secret, err = client.Logical().ReadWithData(dataPath, map[string][]string{
				"version": {fmt.Sprintf("%d", ver)},
			})
		}
	}

	if err != nil {
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "no such host") {
			logging.Errorf("Vault service unavailable: %v", err)
			return types.ErrVaultUnavailable
		}
		if strings.Contains(err.Error(), "permission denied") || strings.Contains(err.Error(), "unauthorized") {
			logging.Errorf("Invalid Vault token or unauthorized: %v", err)
			return types.ErrVaultInvalidAuth
		}
		if strings.Contains(err.Error(), "server is sealed") {
			logging.Errorf("Vault is sealed: %v", err)
			return types.ErrVaultSealed
		}

		logging.Errorf("Secret read failed: %v", err)
		return types.ErrReadSecret
	}

	if secret == nil {
		logging.Errorf("Secret read returned nil")
		return types.ErrReadSecret
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		logging.Errorf("Secret format invalid")
		return types.ErrExtractData
	}

	return outputData(data, types.Quiet)
}
