// vaultreader
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp : 2025/06/04 14:55
// Original filename : src/kv/outputformatter.go

package kv

import (
	"encoding/json"
	"fmt"
	"os"
	"vaultreader/logging"
	"vaultreader/types"
)

func outputData(data map[string]interface{}, suppress bool) int {
	if types.KVSecretField != "" {
		logging.Debugf("Reading field: %s", types.KVSecretField)
		val, found := data[types.KVSecretField]
		if !found {
			logging.Errorf("Field not found: %s", types.KVSecretField)
			return types.ErrFieldNotFound
		}
		if suppress {
			return 0
		}
		if types.OutputFormat == "json" {
			out := map[string]interface{}{types.KVSecretField: val}
			json.NewEncoder(os.Stdout).Encode(out)
		} else {
			fmt.Printf("%v\n", val)
		}
		return 0
	}

	if suppress {
		return 0
	}

	if types.OutputFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(data); err != nil {
			logging.Errorf("JSON encoding failed: %v", err)
			return types.ErrExtractData
		}
	} else {
		for k, v := range data {
			fmt.Printf("%s: %v\n", k, v)
		}
	}
	return 0
}
