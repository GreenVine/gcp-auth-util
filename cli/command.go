package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	c "github.com/GreenVine/gcp-auth-util/common"
	"github.com/go-playground/validator/v10"
	"golang.org/x/oauth2/google"
)

// CommandLineOptions describes arguments accepted by command line interface
type CommandLineOptions struct {
	Audience       string        `validate:"omitempty"`
	AuthSource     string        `validate:"omitempty,oneof='service-account'"`
	Credentials    string        `validate:"base64|file,required"`
	Expires        time.Duration `validate:"omitempty,gt=0"`
	Scopes         stringSlice   `validate:"omitempty,unique,dive,url"`
	TargetAudience string        `validate:"omitempty"`
	Timeout        time.Duration `validate:"omitempty,gt=0"`
	TokenType      string        `validate:"oneof=id access,required"`
	TokenURL       string        `validate:"omitempty,url"`
}

// IsIDTokenType checks if CLI options are representing an ID token request. Otherwise it's an access token.
func (options *CommandLineOptions) IsIDTokenType() bool {
	return options.TokenType == "id"
}

// ParseCommand retrieves arguments from command line
func ParseCommand() (CommandLineOptions, error) {
	options := CommandLineOptions{}

	// Get variables from command line and fallback to default value to environment variable / hardcoded value
	flag.StringVar(
		&options.Audience, "audience",
		c.GetEnvString("AUDIENCE", google.JWTTokenURL), "Token audience")
	flag.StringVar(
		&options.AuthSource, "auth-source",
		c.GetEnvString("AUTH_SOURCE", "service-account"), "Authentication source: service-account")
	flag.StringVar(
		&options.Credentials, "credentials",
		c.GetEnvString("CREDENTIALS", ""), "Base64-encoded credentials or path to the file")
	flag.DurationVar(&options.Expires, "expires",
		c.GetEnvDuration("TOKEN_EXPIRES", c.DefaultTokenExpiry), "Token expiry (in seconds)")
	flag.Var(&options.Scopes, "scopes", "One or more scopes")
	flag.StringVar(
		&options.TargetAudience, "target-audience",
		c.GetEnvString("TARGET_AUDIENCE", ""), "Target audience")
	flag.DurationVar(&options.Timeout, "timeout",
		c.GetEnvDuration("TIMEOUT", c.DefaultRequestContextTimeout), "Operation timeout (in seconds)")
	flag.StringVar(
		&options.TokenType, "token-type",
		c.GetEnvString("TOKEN_TYPE", ""), "Token type: id, access")
	flag.StringVar(
		&options.TokenURL, "token-url",
		c.GetEnvString("TOKEN_URL", google.JWTTokenURL), "OAuth2 authentication endpoint")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "GCP Auth Utility (%s)\n\nUsage of %s:\n", c.BuildTag, os.Args[0])

		flag.PrintDefaults()
	}

	flag.Parse()

	// Validate command line options
	var errFields []string
	if err := validator.New().Struct(options); err != nil {
		errors := err.(validator.ValidationErrors)

		for i, err := range errors {
			errFields = append(errFields, fmt.Sprintf("#%d %s: %s", i+1, err.Namespace(), err.Tag()))
		}
	}

	if len(errFields) > 0 {
		return options, fmt.Errorf("%d argument validation error(s):\n%s",
			len(errFields), strings.Join(errFields, "\r\n"))
	}
	return options, nil
}
