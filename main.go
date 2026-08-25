package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/zalando/go-keyring"
)

// KeyringProvider defines the interface for interacting with the system keyring.
type KeyringProvider interface {
	Get(service, user string) (string, error)
}

// SystemKeyring is the real implementation using go-keyring.
type SystemKeyring struct{}

func (s SystemKeyring) Get(service, user string) (string, error) {
	return keyring.Get(service, user)
}

// keyringProvider is a package-level variable to allow dependency injection during testing.
var keyringProvider KeyringProvider = SystemKeyring{}

func getAuthToken() (string, error) {
	// 1. Check environment variables first
	if token := os.Getenv("GH_TOKEN"); token != "" {
		return token, nil
	}

	// 2. Attempt keyring lookup
	token, err := keyringProvider.Get("gh-cli", "auth-token")
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return "", nil // Safe to proceed unauthenticated
		}
		return "", fmt.Errorf("failed to access system keyring: %w", err) // Abort on system error
	}

	return token, nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	token, err := getAuthToken()
	if err != nil {
		return err
	}

	if token == "" {
		fmt.Println("Running as unauthenticated user")
	} else {
		fmt.Println("Authenticated successfully")
	}
	return nil
}