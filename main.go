package main

import (
	"errors"
	"fmt"
	"os"
)

// Mocking keyring errors
var ErrNotFound = errors.New("secret not found")

func getTokenFromKeyring() (string, error) {
	// Simulate a system error (e.g., dbus connection failure)
	// In a real scenario, this would call the keyring library
	return "", errors.New("failed to access system keyring: dbus connection refused")
}

func getAuthToken() (string, error) {
	// 1. Check environment variables first
	if token := os.Getenv("GH_TOKEN"); token != "" {
		return token, nil
	}

	// 2. Attempt keyring lookup
	token, err := getTokenFromKeyring()
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return "", nil // Safe to proceed unauthenticated
		}
		return "", err // Abort on system error
	}

	return token, nil
}

func main() {
	token, err := getAuthToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if token == "" {
		fmt.Println("Running as unauthenticated user")
	} else {
		fmt.Println("Authenticated successfully")
	}
}