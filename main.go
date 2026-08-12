package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/zalando/go-keyring"
)

const (
	keyringService  = "gh"
	keyringUsername = "oauth_token"
)

// keyringGet is a variable so token resolution can be tested without depending
// on a desktop keyring being available in the test environment.
var keyringGet = keyring.Get

func getTokenFromKeyring() (string, error) {
	return keyringGet(keyringService, keyringUsername)
}

func getAuthToken() (string, error) {
	// Explicit environment tokens bypass the keyring entirely.
	for _, name := range []string{"GH_TOKEN", "GITHUB_TOKEN"} {
		if token := os.Getenv(name); token != "" {
			return token, nil
		}
	}

	token, err := getTokenFromKeyring()
	if err == nil {
		return token, nil
	}
	if errors.Is(err, keyring.ErrNotFound) {
		// A missing credential is different from an inaccessible keyring and may
		// legitimately result in an unauthenticated request.
		return "", nil
	}

	return "", fmt.Errorf("failed to access system keyring: %w", err)
}

func main() {
	token, err := getAuthToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if token == "" {
		fmt.Println("Running as unauthenticated user")
		return
	}
	fmt.Println("Authenticated successfully")
}
