// Hashicorp Vault Client (vaultreader)
// src/cmd/root.go

package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"vaultreader/kv"
	"vaultreader/types"

	hftfx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vaultreader KV_engine secret_path [-f field]",
	Short: "Read-only Vault client for KV v2 secrets",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if kvreadErr := kv.ReadSecrets(args[0], args[1]); kvreadErr != nil {
			fmt.Println(hftfx.SkullBonesSign(kvreadErr.Error()))
			os.Exit(1)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the software version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(hftfx.White(fmt.Sprintf("2.00.00 (2026.04.25), Go version  : " +
			strings.TrimPrefix(runtime.Version(), "go"))))
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
