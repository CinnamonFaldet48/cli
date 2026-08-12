package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/zalando/go-keyring"
)

func withKeyringGet(t *testing.T, get func(string, string) (string, error)) {
	t.Helper()
	original := keyringGet
	keyringGet = get
	t.Cleanup(func() { keyringGet = original })
}

func TestGetAuthTokenUsesGHTokenWithoutKeyring(t *testing.T) {
	t.Setenv("GH_TOKEN", "gh-token")
	t.Setenv("GITHUB_TOKEN", "github-token")
	withKeyringGet(t, func(string, string) (string, error) {
		t.Fatal("keyring should not be queried when an environment token is set")
		return "", nil
	})

	token, err := getAuthToken()
	if err != nil || token != "gh-token" {
		t.Fatalf("getAuthToken() = %q, %v; want %q, nil", token, err, "gh-token")
	}
}

func TestGetAuthTokenUsesGITHUBTokenWithoutKeyring(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "github-token")
	withKeyringGet(t, func(string, string) (string, error) {
		t.Fatal("keyring should not be queried when an environment token is set")
		return "", nil
	})

	token, err := getAuthToken()
	if err != nil || token != "github-token" {
		t.Fatalf("getAuthToken() = %q, %v; want %q, nil", token, err, "github-token")
	}
}

func TestGetAuthTokenAllowsMissingKeyringCredential(t *testing.T) {
	withKeyringGet(t, func(string, string) (string, error) {
		return "", keyring.ErrNotFound
	})

	token, err := getAuthToken()
	if err != nil || token != "" {
		t.Fatalf("getAuthToken() = %q, %v; want empty token and nil error", token, err)
	}
}

func TestGetAuthTokenReturnsKeyringFailure(t *testing.T) {
	withKeyringGet(t, func(string, string) (string, error) {
		return "", errors.New("keyring locked")
	})

	token, err := getAuthToken()
	if token != "" {
		t.Fatalf("getAuthToken() returned token %q; want empty token", token)
	}
	if err == nil || !strings.Contains(err.Error(), "failed to access system keyring: keyring locked") {
		t.Fatalf("getAuthToken() error = %v; want wrapped keyring failure", err)
	}
}
