// vaultreader
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp : 2025/06/04 10:04
// Original filename : src/kv/helpers.go

package kv

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"os"
	"sort"
	"strconv"
	"strings"
	"vaultreader/logging"
	"vaultreader/types"

	hfl "github.com/jeanfrancoisgratton/helperFunctions/v3/logging"
	hftfx "github.com/jeanfrancoisgratton/helperFunctions/v3/terminalfx"
)

func findLatestAvailableVersion(client *api.Client, metaPath string) (int, error) {
	meta, err := client.Logical().Read(metaPath)
	if err != nil || meta == nil {
		return 0, fmt.Errorf("unable to fetch metadata")
	}

	rawVersions, ok := meta.Data["versions"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("versions metadata not found")
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

func setGlobals() int {
	if types.VaultAuthToken == "" {
		logging.Debugf("No vault auth token were provided at command line")
		types.VaultAuthToken = os.Getenv("VAULT_TOKEN")
	}
	if types.VaultAuthToken == "" {
		logging.Infof("$VAULT_TOKEN environment variable not set")
		homeDir, err := os.UserHomeDir()
		if err == nil {
			data, err := os.ReadFile(homeDir + "/.vault-token")
			if err == nil {
				types.VaultAuthToken = strings.TrimSpace(string(data))
			}
		}
	}
	if types.VaultAuthToken == "" {
		logging.Errorf("Vault token missing")
		return types.ErrNoToken
	}

	if types.VaultServerAddress == "" {
		logging.Debugf("No vault server address were provided at command line")
		types.VaultServerAddress = os.Getenv("VAULT_ADDR")
	}
	if types.VaultServerAddress == "" {
		logging.Errorf("Vault address missing")
		return types.ErrNoAddress
	}
	return 0
}
