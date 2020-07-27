# gcp-auth-util

This is a utility library that requests short-lived tokens for [GCP](https://cloud.google.com) [Service Accounts](https://cloud.google.com/iam/docs/service-accounts). Command line version is also available as a standalone utility tool.

## Installation

`go get -u github.com/GreenVine/gcp-auth-util/auth`

## Usage

### As a Library

Build a `TokenRequestOptions` with appropriate settings, then retrieve the token by authenticating with service account credentials.

To request an ID token, provide `TargetAudience` and set `UseIDToken` to `true`:

```go
import "github.com/GreenVine/gcp-auth-util/auth"

options := &auth.TokenRequestOptions{
    Audience:       "https://oauth2.googleapis.com/token", // Audience for the token
    Expires:        3600, // Token expiry (in seconds)
    TargetAudience: "https://myapp.example.com", // Target audience for the token
    TokenURL:       "https://oauth2.googleapis.com/token", // OAuth2 endpoint
    UseIDToken:     true, // true: request an ID token; false: request an access token
}
```

To request an access token instead, provide at least one scope and set `UseIDToken` to `false`:

```go
options := &auth.TokenRequestOptions{
    // ...
    Scopes:         []string{ // List of scopes
        "https://www.googleapis.com/auth/userinfo.email",
    },
    // ...
    UseIDToken:     false,
}
```

Then authenticate by a service account and request a token:

```go
// Pass credentials via JSON file
token, err := auth.GetTokenByServiceAccountFile(context.Background(), "/path/to/service-account.json", options)

// Pass credentials via Base64-encoded contents
token, err := auth.GetTokenByServiceAccountFile(context.Background(), []byte("base64-encoded-json-file-contents"), options)

// Print the token
fmt.Println(token.AccessToken)
```

### As a Command Line Tool

Build the project with: `go build -o gau -ldflags="-s -w"`, where `gau` is the resulting binary.

#### CLI Options

```
-audience string
    Token audience (default "https://oauth2.googleapis.com/token")
-auth-source string
    Authentication source: service-account (default "service-account")
-credentials string
    Base64-encoded credentials or path to the file
-expires duration
    Token expiry (in seconds) (default 1h0m0s)
-scopes value
    One or more scopes
-target-audience string
    Target audience
-timeout duration
    Operation timeout (in seconds) (default 15s)
-token-type string
    Token type: id, access
-token-url string
    OAuth2 authentication endpoint (default "https://oauth2.googleapis.com/token")
```

#### CLI Examples

1. Request a 5-min-long ID token, authenticate by JSON file:

    `gau -credentials /path/to/service-account.json -token-type id -expires 300 -target-audience "https://myapp.example.com"`

2. Request an access token, authenticate by Base64-encoded credentials:

    `gau -credentials "base64-encoded" -token-type access -scopes https://www.googleapis.com/auth/userinfo.email -scopes https://example.com/another-scope`

## License
GNU General Public License v3.0 or later

See [COPYING](COPYING) to see the full text.
