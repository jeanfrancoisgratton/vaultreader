# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html) since 2.2.0.

## [2.5.0] - 2026-11-07

### Fixed
- Reading a non-existent field from an existing secret now exits with `ErrFieldNotFound` (6) instead of `ErrReadSecret` (4)
- Reading a field (`-f`) from a non-existent secret, or a non-existent/deleted secret version (`-v`), now exits with `ErrInvalidPath` (7) instead of `ErrReadSecret` (4)

### Removed
- Unreachable single-field and JSON code paths in the whole-secret text formatter

### Changed
- Minor code refactoring: the data types what were needlessly sitting in their own types/ packages are now part of the kv/ package

## [2.4.0] - 2026-09-22

### Added
- `--file` flag to write a secret straight to a file instead of stdout
- Windows (.msi) packaging support
- Test coverage

### Changed
- Packaging modernized, with dynamic version numbering

## [2.3.1] - 2026-08-11

### Removed
- Shell completion scripts

## [2.3.0] - 2026-08-01

### Changed
- Updated the customError package: errors are now sent to stderr

## [2.2.0] - 2026-06-24

### Changed
- Naming scheme now aligns with SemVer

### Fixed
- Missing dependencies at post-install

## [2.01.00] - 2026-06-03

### Changed
- Go upgraded to v1.26.4
- Now uses the refactored vaultLib package

## [2.00.01] - 2026-05-13

### Added
- Arch Linux support

### Changed
- Refactored binary packaging

### Removed
- Logging support
- All pre/post install scripts

## [2.00.00] - 2026-04-25

### Changed
- The tool is no longer self-contained: it now calls vaultLib for Vault API calls

## [1.41.00] - 2026-01-19

### Changed
- Quieter output, useful when called from a script

## [1.40.04] - 2025-11-28

### Changed
- Migrated to helperFunctions/v4

## [1.40.02] - 2025-11-14

### Fixed
- Output when `-o json` is used

## [1.40.00] - 2025-11-11

### Changed
- Cleaned up error handling (phase I)
- Build script enhancements

### Fixed
- Wrong environment variable name (`VAULT_ADDR`)

## [1.30.00] - 2025-11-06

### Added
- Extended logging options as per helperFunctions v3.06+

### Changed
- Go version bump
- Prettified output
- Migrated logging to the helperFunctions package
- Moved logging to `/var/log/vaultreader.log`
- Cleaned up Alpine packaging

## [1.23.00] - 2025-11-02

### Added
- Tab completion

### Changed
- Go version bump
- Binary stripping at build stage

## [1.22.01] - 2025-09-11

### Changed
- Go version bump

## [1.22.00] - 2025-07-25

### Changed
- Go version bump
- Moved logfile to `$HOME/.local/state/vaultreader/`, removing the need to run as root

## [1.21.00] - 2025-07-01

### Changed
- More error handling

## [1.20.00] - 2025-07-01

### Added
- Error when the secret path is missing

## [1.10.01] - 2025-06-05

### Fixed
- Mistyped `VAULT_ADDR` variable

## [1.10.00] - 2025-06-04

### Added
- Logging capabilities

## [1.01.00] - 2025-06-04

### Changed
- Better secret version handling

## [1.00.00] - 2025-05-30

### Added
- Initial version
