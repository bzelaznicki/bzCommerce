package main

import (
	"html"
	"regexp"
	"strings"
	"unicode"
)

// sanitizeInput removes potentially dangerous characters and HTML from user input
func sanitizeInput(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)
	
	// Escape HTML entities to prevent XSS
	input = html.EscapeString(input)
	
	return input
}

// sanitizeName sanitizes full name input
func sanitizeName(name string) string {
	// Trim whitespace first
	name = strings.TrimSpace(name)
	
	// Allow only letters, spaces, hyphens, and apostrophes
	reg := regexp.MustCompile(`[^a-zA-Z\s\-']+`)
	name = reg.ReplaceAllString(name, "")
	
	// Remove excessive whitespace
	reg = regexp.MustCompile(`\s+`)
	name = reg.ReplaceAllString(name, " ")
	
	return strings.TrimSpace(name)
}

// validatePasswordStrength checks if password meets security requirements
func validatePasswordStrength(password string) bool {
	if len(password) < MinPasswordLength {
		return false
	}
	
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)
	
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	
	// Require at least 3 of the 4 character types for passwords >= 8 chars
	criteriaCount := 0
	if hasUpper { criteriaCount++ }
	if hasLower { criteriaCount++ }
	if hasNumber { criteriaCount++ }
	if hasSpecial { criteriaCount++ }
	
	return criteriaCount >= 3
}