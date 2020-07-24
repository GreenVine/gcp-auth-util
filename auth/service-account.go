package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
	"time"

	c "github.com/GreenVine/gcp-auth-util/common"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
)

// TokenRequestOptions is a subset of `jwt.Config` that can be exposed to the user
type TokenRequestOptions struct {
	Audience       string        `json:"audience" validate:"omitempty"`
	Email          string        `json:"client_email" validate:"omitempty,email"`
	Expires        time.Duration `json:"expires" validate:"gt=0"`
	PrivateKey     string        `json:"private_key" validate:"required"`
	PrivateKeyID   string        `json:"private_key_id" validate:"omitempty"`
	Scopes         []string      `json:"scopes" validate:"required_without=UseIDToken,min=1,dive,url"`
	Subject        string        `json:"subject" validate:"omitempty"`
	TargetAudience string        `json:"target_audience" validate:"required_with=UseIDToken,url|min=0"`
	TokenURL       string        `json:"token_url" validate:"required,url"`

	UseIDToken bool `json:"use_id_token" validate:"omitempty"` // internal field (to be overwritten)
}

// GetTokenByServiceAccountFile retrieves OAuth2 token from a service account credential file
func GetTokenByServiceAccountFile(ctx context.Context, serviceAccountFile string, options *TokenRequestOptions) (*oauth2.Token, error) {
	content, err := c.ReadEntireFile(serviceAccountFile)
	if err != nil {
		return nil, err
	}

	return GetTokenByServiceAccount(ctx, content, options)
}

// GetTokenByServiceAccount retrieves OAuth2 token from a service account JSON string
func GetTokenByServiceAccount(ctx context.Context, serviceAccount []byte, options *TokenRequestOptions) (*oauth2.Token, error) {
	if options == nil {
		return nil, fmt.Errorf("empty token request options")
	}
	if err := json.Unmarshal(serviceAccount, options); err != nil {
		return nil, err
	}

	return getToken(ctx, options)
}

// getToken requests an ID or access token from OAuth2 endpoint
func getToken(ctx context.Context, options *TokenRequestOptions) (*oauth2.Token, error) {
	// Validate options
	var errFields []string

	if err := validator.New().Struct(options); err != nil {
		errors := err.(validator.ValidationErrors)

		for i, err := range errors {
			errFields = append(errFields, fmt.Sprintf("#%d %s: %s", i+1, err.Namespace(), err.Tag()))
		}
	}

	// Has one or more validation errors
	if len(errFields) > 0 {
		return nil, fmt.Errorf("%d validation error(s) in token request options:\n%s",
			len(errFields), strings.Join(errFields, "\r\n"))
	}

	// Populate user-supplied information into JWT config
	tokenConfig := &jwt.Config{
		Audience:     options.Audience,
		Email:        options.Email,
		Expires:      options.Expires,
		PrivateKey:   []byte(options.PrivateKey),
		PrivateKeyID: options.PrivateKeyID,
		Subject:      options.Subject,
		TokenURL:     options.TokenURL,
		UseIDToken:   options.UseIDToken,
	}

	// Populate token type specific fields
	if options.UseIDToken {
		tokenConfig.PrivateClaims = map[string]interface{}{
			"target_audience": options.TargetAudience,
		}
	} else {
		tokenConfig.Scopes = options.Scopes
	}

	// Generate token
	return tokenConfig.TokenSource(ctx).Token()
}
