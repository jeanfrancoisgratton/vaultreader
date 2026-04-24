// Hashicorp Vault Client (vaultreader)
// src/cmd/root.go

package cmd

import (
	"fmt"
	"os"
	"runtime"

	"vaultreader/kv"
	"vaultreader/types"

	hfl "github.com/jeanfrancoisgratton/helperFunctions/v5/logging"
	hftfx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vaultreader -m KV_engine secret_path -f [field]",
	Short: "Read-only Vault client for KV v2 secrets",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var xcode = 0
		if kvreadErr := kv.ReadSecrets(args[0]); kvreadErr != nil {
			xcode = kvreadErr.Code
		}

		if types.LogLevel != "none" {
			hfl.Close()
		}
		os.Exit(xcode)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the software version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(hftfx.White(fmt.Sprintf("1.41.00-%s (2026.01.19)", runtime.GOARCH)))
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(10)
	}
}

func init() {
	rootCmd.DisableAutoGenTag = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(completionCmd, versionCmd)

	rootCmd.PersistentFlags().StringVarP(&types.VaultAuthToken, "token", "t", "", "Vault token (or use VAULT_TOKEN)")
	rootCmd.PersistentFlags().StringVarP(&types.VaultServerAddress, "vaultaddress", "a", types.VaultServerAddress, "Vault server address (or use VAULT_ADDRESS)")
	rootCmd.PersistentFlags().StringVarP(&types.KVEngineMountPath, "mount", "m", "", "KV v2 mount path (required)")
	rootCmd.PersistentFlags().IntVarP(&types.KVSecretVersion, "version", "v", 0, "Secret version (0 = latest available)")
	rootCmd.PersistentFlags().StringVarP(&types.KVSecretField, "field", "f", "", "Specific field to display")
	rootCmd.PersistentFlags().BoolVarP(&types.Quiet, "quiet", "q", false, "Suppress all stdout output")
	rootCmd.PersistentFlags().StringVarP(&types.OutputFormat, "output", "o", "text", "Output format: text|json")
	//rootCmd.MarkPersistentFlagRequired("mount")
}
