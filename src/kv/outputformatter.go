// vaultreader
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp : 2025/06/04 14:55
// Original filename : src/kv/outputformatter.go

package kv

import (
	"fmt"
)

// outputData prints a whole secret as "key: value" lines. Single-field reads
// and JSON output are handled upstream in read.go and never reach here.
func outputData(data map[string]interface{}, suppress bool) {
	if suppress {
		return
	}
	for k, v := range data {
		fmt.Printf("%s: %v\n", k, v)
	}
}
