| Release | Date        | Comments                                                                                                               |
|---------|-------------|------------------------------------------------------------------------------------------------------------------------|
| 2.00.01 | 2026.05.13  | Refactored binary packaging, adding archlinux support                                                                  |
| 2.00.00 | 2026.04.25  | The tool is no longer self-contained, it now calls vaultLib for Vault API calls                                        |
| 1.41.00 | 2026.01.19  | Output is quieter, usefull for when called from a script                                                               |
| 1.40.04 | 2025.11.28  | Migrated to helperFunctions/v4                                                                                         |
| 1.40.02 | 2025.11.14  | Fixed output when -o == json                                                                                           |
| 1.40.00 | 2025.11.11  | Fixed wrong environment variable name (VAULT_ADDR)<br>Cleaned up error handling (phase I)<br>Build script enhancements |
| 1.30.00 | 2025.11.06  | GO version bump<br>Prettyfied output<br>Migrated logging to the helperFunctions package<br>Cleaned Alpine packaging    |
| 1.30.00 | _continued_ | Moved logging to /var/log/vaultreader.log<br>Extended logging options as per helperFunctions v3.06+                    |
| 1.23.00 | 2025.11.02  | GO version bump<br>Tab completion<br>Binary striping at build-stage                                                    |
| 1.22.01 | 2025.09.11  | GO version bump                                                                                                        |
| 1.22.00 | 2025.07.25  | GO version bump<br>Moved logfile to `$HOME/.local/state/vaultreader/`,<br>thus removing the need of running under root |
| 1.21.00 | 2025.07.01  | More error handling                                                                                                    |
| 1.20.00 | 2025.07.01  | Added error when secret path is miising                                                                                |
| 1.10.01 | 2025.06.05  | Fixed mistyped VAULT_ADDR var                                                                                          |
| 1.10.00 | 2025.06.04  | Added logging capabilities                                                                                             |
| 1.01.00 | 2025.06.04  | Better secret version handling                                                                                         |
| 1.00.00 | 2025.05.30  | Initial version                                                                                                        |