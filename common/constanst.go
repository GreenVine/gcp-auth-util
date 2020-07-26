package common

import "time"

// Application defaults
const DefaultTokenExpiry = time.Hour
const DefaultRequestContextTimeout = 15 * time.Second

// Application exit codes
const GenericErrorExitCode = 1
const InvalidArgumentErrorExitCode = 2
const TokenErrorExitCode = 3
