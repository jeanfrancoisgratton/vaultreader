// vclt
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/05/29 16:02
// Original filename: src/types/types.go

package types

var (
	VaultAuthToken     string
	VaultServerAddress string
	KVEngineMountPath  string
	KVSecretVersion    int
	KVSecretField      string
	Quiet              bool
	OutputFormat       string
	LogLevel           string
)

// Exit codes
const (
	ErrNoToken = iota + 1
	ErrNoAddress
	ErrVaultInit
	ErrVaultAuthTokenMissing
	ErrVaultServerAddressMissing
	ErrReadSecret
	ErrExtractData
	ErrFieldNotFound
	ErrInvalidPath
	ErrVaultUnavailable
	ErrVaultSealed
	ErrVaultInvalidAuth
)

var ErrorMessages = map[int]string{
	ErrNoToken:          "No Vault auth token provided",
	ErrNoAddress:        "No Vault server address provided",
	ErrVaultInit:        "Error initializing Vault client",
	ErrReadSecret:       "Error reading secret from Vault",
	ErrExtractData:      "Error extracting secret data",
	ErrFieldNotFound:    "Requested field not found in secret",
	ErrInvalidPath:      "Secret path does not exist",
	ErrVaultUnavailable: "Vault server unavailable",
	ErrVaultSealed:      "Vault is sealed",
	ErrVaultInvalidAuth: "Vault auth token is invalid",
}
