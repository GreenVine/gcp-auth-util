package main

import (
	"context"
	"fmt"
	"github.com/GreenVine/gcp-auth-util/auth"
	"github.com/GreenVine/gcp-auth-util/cli"
	c "github.com/GreenVine/gcp-auth-util/common"
	"golang.org/x/oauth2"
	"log"
	"strings"
)

func main() {
	// Define shared variables
	var token *oauth2.Token
	var err error

	// Parse command line arguments
	args, err := cli.ParseCommand()
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Build request options
	requestContext, cancel := context.WithTimeout(context.Background(), args.Timeout)
	defer cancel()

	requestOptions := &auth.TokenRequestOptions{
		Audience:       args.Audience,
		Expires:        args.Expires,
		Scopes:         args.Scopes,
		TargetAudience: args.TargetAudience,
		TokenURL:       args.TokenURL,
		UseIDToken:     args.IsIDTokenType(),
	}

	// Get token for each authentication source
	switch args.AuthSource {
	case "service-account":
		if c.IsFileExists(args.Credentials) { // treat credential as a file path
			token, err = auth.GetTokenByServiceAccountFile(requestContext, args.Credentials, requestOptions)
		} else {
			token, err = auth.GetTokenByServiceAccount(requestContext, []byte(args.Credentials), requestOptions)
		}
	default:
		log.Panicf("Unexpected authentication source: %v", args.AuthSource)
	}

	if token != nil && err == nil {
		accessToken := strings.TrimSpace(token.AccessToken)

		if strings.TrimSpace(accessToken) != "" { // successfully retrieved the token
			fmt.Printf("%s", accessToken)
		} else {
			log.Fatalf("Empty token returned from the server")
		}
	} else {
		log.Fatalf("Error occurred: %v", err)
	}
}
