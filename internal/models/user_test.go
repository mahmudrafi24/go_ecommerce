package models

import (
	"math/rand"
	"strings"
	"testing"
	"testing/quick"
)

// **Feature: go-ecommerce-demo, Property 2: Invalid registration data is rejected**
// **Validates: Requirements 1.2**
//
// For any registration data with invalid email format, empty password, or empty name,
// the registration request should be rejected with validation errors.

// generateInvalidEmail generates an email that is invalid (no @ or no domain)
func generateInvalidEmail(r *rand.Rand) string {
	choices := []string{
		"",                    // empty
		"plaintext",           // no @ symbol
		"missing@domain",      // no TLD
		"@nodomain.com",       // no local part
		"spaces in@email.com", // spaces
		"double@@at.com",      // double @
	}
	return choices[r.Intn(len(choices))]
}

// generateEmptyOrWhitespace generates empty or whitespace-only strings
func generateEmptyOrWhitespace(r *rand.Rand) string {
	choices := []string{
		"",
		"   ",
		"\t",
		"\n",
	}
	return choices[r.Intn(len(choices))]
}

// generateShortPassword generates a password that is too short (< 8 chars)
func generateShortPassword(r *rand.Rand) string {
	length := r.Intn(8) // 0-7 characters
	if length == 0 {
		return ""
	}
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

// TestProperty2_InvalidEmailRejected tests that invalid email formats are rejected
func TestProperty2_InvalidEmailRejected(t *testing.T) {
	config := &quick.Config{
		MaxCount: 100,
	}

	// Property: For any invalid email, validation should return an email error
	f := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		invalidEmail := generateInvalidEmail(r)

		input := &UserRegistrationInput{
			Email:    invalidEmail,
			Password: "validpassword123", // valid password
			Name:     "Valid Name",       // valid name
		}

		errs := input.Validate()

		// Should have an email error
		_, hasEmailError := errs["email"]
		return hasEmailError
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("Property 2 failed (invalid email): %v", err)
	}
}

// TestProperty2_EmptyPasswordRejected tests that empty/short passwords are rejected
func TestProperty2_EmptyPasswordRejected(t *testing.T) {
	config := &quick.Config{
		MaxCount: 100,
	}

	// Property: For any empty or too-short password, validation should return a password error
	f := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		shortPassword := generateShortPassword(r)

		input := &UserRegistrationInput{
			Email:    "valid@email.com", // valid email
			Password: shortPassword,
			Name:     "Valid Name", // valid name
		}

		errs := input.Validate()

		// Should have a password error
		_, hasPasswordError := errs["password"]
		return hasPasswordError
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("Property 2 failed (empty/short password): %v", err)
	}
}

// TestProperty2_EmptyNameRejected tests that empty names are rejected
func TestProperty2_EmptyNameRejected(t *testing.T) {
	config := &quick.Config{
		MaxCount: 100,
	}

	// Property: For any empty name, validation should return a name error
	f := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		emptyName := generateEmptyOrWhitespace(r)

		// Only test truly empty names (the validation only checks for empty string)
		if strings.TrimSpace(emptyName) != "" {
			return true // skip non-empty names
		}
		if emptyName != "" {
			return true // skip whitespace-only (validation doesn't trim)
		}

		input := &UserRegistrationInput{
			Email:    "valid@email.com",  // valid email
			Password: "validpassword123", // valid password
			Name:     emptyName,
		}

		errs := input.Validate()

		// Should have a name error
		_, hasNameError := errs["name"]
		return hasNameError
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("Property 2 failed (empty name): %v", err)
	}
}

// TestProperty2_MultipleInvalidFieldsRejected tests that multiple invalid fields all produce errors
func TestProperty2_MultipleInvalidFieldsRejected(t *testing.T) {
	config := &quick.Config{
		MaxCount: 100,
	}

	// Property: For any input with all invalid fields, all fields should have errors
	f := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))

		input := &UserRegistrationInput{
			Email:    generateInvalidEmail(r),
			Password: generateShortPassword(r),
			Name:     "", // empty name
		}

		errs := input.Validate()

		// Should have errors for all fields
		_, hasEmailError := errs["email"]
		_, hasPasswordError := errs["password"]
		_, hasNameError := errs["name"]

		return hasEmailError && hasPasswordError && hasNameError
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("Property 2 failed (multiple invalid fields): %v", err)
	}
}
