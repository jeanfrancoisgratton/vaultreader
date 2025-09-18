// Hashicorp Vault Client (vaultreader)
// src/cmd/root.go

package cmd

import (
	"fmt"
	"os"
	"vaultreader/kv"
	"vaultreader/types"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vaultreader [secret path]",
	Short: "Read-only Vault client for KV v2 secrets",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		exitCode := kv.ReadSecrets(args[0])
		if !types.Quiet && exitCode != 0 {
			if msg, ok := types.ErrorMessages[exitCode]; ok {
				fmt.Fprintln(os.Stderr, msg)
			}
		}
		os.Exit(exitCode)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the software version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vaultreader version 1.22.01 (2025.09.11)")
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
	rootCmd.AddCommand(versionCmd)
	unseal := unsealCmd()

	// Show only when running as root; otherwise keep it hidden but usable
	if os.Geteuid() != 0 {
		unseal.Hidden = true
	} else {
		unseal.Hidden = false
	}

	rootCmd.AddCommand(unseal)

	rootCmd.PersistentFlags().StringVarP(&types.VaultAuthToken, "token", "t", "", "Vault token (or use VAULT_TOKEN)")
	rootCmd.PersistentFlags().StringVarP(&types.VaultServerAddress, "vaultaddress", "a", types.VaultServerAddress, "Vault server address (or use VAULT_ADDRESS)")
	rootCmd.PersistentFlags().StringVarP(&types.KVEngineMountPath, "mount", "m", "", "KV v2 mount path (required)")
	rootCmd.PersistentFlags().IntVarP(&types.KVSecretVersion, "version", "v", 0, "Secret version (0 = latest available)")
	rootCmd.PersistentFlags().StringVarP(&types.KVSecretField, "field", "f", "", "Specific field to display")
	rootCmd.PersistentFlags().BoolVarP(&types.Quiet, "quiet", "q", false, "Suppress all stdout output")
	rootCmd.PersistentFlags().StringVarP(&types.OutputFormat, "output", "o", "text", "Output format: text|json")
	rootCmd.PersistentFlags().StringVarP(&types.LogLevel, "loglevel", "l", "error", "Log level: debug|info|error")

	//rootCmd.MarkPersistentFlagRequired("mount")
}
