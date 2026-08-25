package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/zalando/go-keyring"
)

// MockKeyring simulates the keyring for testing purposes.
type MockKeyring struct {
	Token string
	Err   error
}

func (m *MockKeyring) Get(service, user string) (string, error) {
	return m.Token, m.Err
}

func TestGetAuthToken(t *testing.T) {
	// Backup and restore environment variables
	origGH := os.Getenv("GH_TOKEN")
	defer os.Setenv("GH_TOKEN", origGH)

	// Backup and restore the keyring provider
	origProvider := keyringProvider
	defer func() { keyringProvider = origProvider }()

	tests := []struct {
		name          string
		envToken      string
		keyringToken  string
		keyringErr    error
		expectedToken string
		expectedErr   string
	}{
		{
			name:          "env token takes precedence",
			envToken:      "env-token",
			keyringToken:  "keyring-token",
			keyringErr:    nil,
			expectedToken: "env-token",
			expectedErr:   "",
		},
		{
			name:          "keyring token is used if no env token",
			envToken:      "",
			keyringToken:  "keyring-token",
			keyringErr:    nil,
			expectedToken: "keyring-token",
			expectedErr:   "",
		},
		{
			name:          "ErrNotFound falls back silently",
			envToken:      "",
			keyringToken:  "",
			keyringErr:    keyring.ErrNotFound,
			expectedToken: "",
			expectedErr:   "",
		},
		{
			name:          "hard system error fails fast",
			envToken:      "",
			keyringToken:  "",
			keyringErr:    errors.New("dbus connection refused"),
			expectedToken: "",
			expectedErr:   "failed to access system keyring: dbus connection refused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envToken != "" {
				os.Setenv("GH_TOKEN", tt.envToken)
			} else {
				os.Unsetenv("GH_TOKEN")
			}

			keyringProvider = &MockKeyring{
				Token: tt.keyringToken,
				Err:   tt.keyringErr,
			}

			token, err := getAuthToken()
			if err != nil {
				if tt.expectedErr == "" {
					t.Fatalf("unexpected error: %v", err)
				}
				if !strings.Contains(err.Error(), tt.expectedErr) {
					t.Errorf("expected error containing %q, got %q", tt.expectedErr, err.Error())
				}
			} else {
				if tt.expectedErr != "" {
					t.Fatalf("expected error containing %q, got nil", tt.expectedErr)
				}
			}

			if token != tt.expectedToken {
				t.Errorf("expected token %q, got %q", tt.expectedToken, token)
			}
		})
	}
}

// TestMainExitCode tests the actual exit code of the application when a hard system error occurs.
func TestMainExitCode(t *testing.T) {
	if os.Getenv("CRASH_TEST") == "1" {
		// Set up the mock to return a hard error
		keyringProvider = &MockKeyring{
			Err: errors.New("simulated hard dbus error"),
		}
		os.Unsetenv("GH_TOKEN")
		
		// Run main, it should os.Exit(1)
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMainExitCode")
	cmd.Env = append(os.Environ(), "CRASH_TEST=1")
	out, err := cmd.CombinedOutput()
	
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		// Successfully failed with a non-zero exit code
		outStr := string(out)
		if !strings.Contains(outStr, "Error: failed to access system keyring: simulated hard dbus error") {
			t.Errorf("expected stderr to contain error message, got %q", outStr)
		}
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1. Output: %s", err, string(out))
}
