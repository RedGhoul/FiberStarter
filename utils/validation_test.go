package utils

import (
	"testing"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		errorMsg string
	}{
		{"valid username", "john_doe", false, ""},
		{"valid short", "abc", false, ""},
		{"too short", "ab", true, "username must be at least 3 characters"},
		{"too long", "abcdefghijklmnopqrstuvwxyz12345", true, "username must be at most 30 characters"},
		{"invalid chars", "john doe", true, "username can only contain letters, numbers, and underscores"},
		{"special chars", "john@doe", true, "username can only contain letters, numbers, and underscores"},
		{"numbers only", "123", false, ""},
		{"underscore", "a_b_c", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if err != nil && tt.errorMsg != "" && err.Error() != tt.errorMsg {
				t.Errorf("ValidateUsername(%q) error = %q, want %q", tt.input, err.Error(), tt.errorMsg)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid email", "user@example.com", false},
		{"valid with subdomain", "user@mail.example.com", false},
		{"empty", "", true},
		{"no at sign", "userexample.com", true},
		{"no domain", "user@", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid password", "password1", false},
		{"valid complex", "MyP@ssw0rd!", false},
		{"too short", "pass1", true},
		{"no digits", "password", true},
		{"no letters", "12345678", true},
		{"exactly 8 chars valid", "abcdefg1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
