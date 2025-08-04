package main

import (
	"testing"
)

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		password string
		expected bool
		name     string
	}{
		{"", false, "empty password"},
		{"short", false, "too short"},
		{"12345678", false, "only numbers"},
		{"abcdefgh", false, "only lowercase"},
		{"ABCDEFGH", false, "only uppercase"},
		{"!@#$%^&*", false, "only special chars"},
		{"Password1", true, "strong password"},
		{"Pass123!", true, "strong password with special"},
		{"myPassword2", true, "strong password"},
		{"Test@123", true, "strong password"},
		{"weakpass", false, "weak password"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validatePasswordStrength(tt.password)
			if result != tt.expected {
				t.Errorf("validatePasswordStrength(%q) = %v, want %v", tt.password, result, tt.expected)
			}
		})
	}
}

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		name     string
	}{
		{"John Doe", "John Doe", "normal name"},
		{"John123Doe", "JohnDoe", "name with numbers"},
		{"John  Doe", "John Doe", "name with extra spaces"},
		{"<script>alert('xss')</script>", "scriptalert'xss'script", "XSS attempt - tags removed"},
		{"O'Connor", "O'Connor", "name with apostrophe"},
		{"Jean-Pierre", "Jean-Pierre", "name with hyphen"},
		{"", "", "empty name"},
		{"   ", "", "whitespace only"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeName(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
		name     string
	}{
		{"test@example.com", true, "valid email"},
		{"user+tag@domain.org", true, "email with plus"},
		{"invalid.email", false, "no @ symbol"},
		{"@domain.com", false, "no username"},
		{"user@", false, "no domain"},
		{"", false, "empty email"},
		{"   ", false, "whitespace only"},
		{"user@domain", true, "email without TLD (technically valid)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidEmail(tt.email)
			if result != tt.expected {
				t.Errorf("isValidEmail(%q) = %v, want %v", tt.email, result, tt.expected)
			}
		})
	}
}