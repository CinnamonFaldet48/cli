package main

import (
	"errors"
	"fmt"
	"github.com/zalando/go-keyring"
)

// GetToken retrieves the token from the keyring, handling errors appropriately.
func GetToken(service, username string) (string, error) {
	token, err := keyring.Get(service, username)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return "", nil
		}
		return "", fmt.Errorf("keyring lookup failed: %w", err)
	}
	return token, nil
}

func main() {
	fmt.Println("Hello, Bounty Hunter!")
}