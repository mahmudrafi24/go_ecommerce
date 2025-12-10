package models

import (
	"errors"
	"regexp"
	"time"
)

// User represents a registered user in the system
type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	Role         string    `json:"role"` // "user" or "admin"
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserRole constants
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// Validation errors
var (
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrEmptyEmail       = errors.New("email is required")
	ErrEmptyPassword    = errors.New("password is required")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrEmptyName        = errors.New("name is required")
	ErrInvalidRole      = errors.New("role must be 'user' or 'admin'")
)

// emailRegex is a simple regex for email validation
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail checks if the email format is valid
func ValidateEmail(email string) error {
	if email == "" {
		return ErrEmptyEmail
	}
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

// ValidatePassword checks if the password meets requirements
func ValidatePassword(password string) error {
	if password == "" {
		return ErrEmptyPassword
	}
	if len(password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

// ValidateName checks if the name is valid
func ValidateName(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	return nil
}

// ValidateRole checks if the role is valid
func ValidateRole(role string) error {
	if role != RoleUser && role != RoleAdmin {
		return ErrInvalidRole
	}
	return nil
}

// ValidationErrors holds multiple validation errors
type ValidationErrors map[string]string

// Error implements the error interface
func (v ValidationErrors) Error() string {
	return "validation failed"
}

// HasErrors returns true if there are validation errors
func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}

// UserRegistrationInput represents the input for user registration
type UserRegistrationInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// Validate validates the registration input
func (u *UserRegistrationInput) Validate() ValidationErrors {
	errs := make(ValidationErrors)

	if err := ValidateEmail(u.Email); err != nil {
		errs["email"] = err.Error()
	}
	if err := ValidatePassword(u.Password); err != nil {
		errs["password"] = err.Error()
	}
	if err := ValidateName(u.Name); err != nil {
		errs["name"] = err.Error()
	}

	return errs
}

// UserLoginInput represents the input for user login
type UserLoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate validates the login input
func (u *UserLoginInput) Validate() ValidationErrors {
	errs := make(ValidationErrors)

	if err := ValidateEmail(u.Email); err != nil {
		errs["email"] = err.Error()
	}
	if u.Password == "" {
		errs["password"] = ErrEmptyPassword.Error()
	}

	return errs
}

// IsAdmin returns true if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
