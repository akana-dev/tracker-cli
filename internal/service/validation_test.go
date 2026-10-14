package service

import (
	"strings"
	"testing"
)

func makeString(length int) string {
	return strings.Repeat("a", length)
}

func TestValidateTitle(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{"empty", "", true},
		{"whitespace only", "   ", true},
		{"too long", makeString(MaxTitleLength + 1), true},
		{"max length", makeString(MaxTitleLength), false},
		{"valid short", "Task title", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTitle(tt.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTitle(%q) error = %v, wantErr %v", tt.title, err, tt.wantErr)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"empty", "", true},
		{"whitespace only", "   ", true},
		{"too long", makeString(MaxUsernameLength + 1), true},
		{"invalid chars", "user@name", true},
		{"valid with dots dashes underscores", "user.name-test_1", false},
		{"valid alphanumeric", "username123", false},
		{"max length", makeString(MaxUsernameLength), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername(%q) error = %v, wantErr %v", tt.username, err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"empty", "", true},
		{"missing at", "example.com", true},
		{"missing domain", "user@", true},
		{"missing tld", "user@example", true},
		{"invalid chars local", "user name@example.com", true},
		{"too long", makeString(MaxEmailLength-10) + "@example.com", true},
		{"valid simple", "user@example.com", false},
		{"valid with dots and plus", "user.name+tag@example.co.uk", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail(%q) error = %v, wantErr %v", tt.email, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"empty", "", true},
		{"too short", "12345", true},
		{"too long", makeString(129), true},
		{"min length", "123456", false},
		{"valid", "secure-password-123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword(%q) error = %v, wantErr %v", tt.password, err, tt.wantErr)
			}
		})
	}
}

func TestValidateCompanyName(t *testing.T) {
	tests := []struct {
		name    string
		company string
		wantErr bool
	}{
		{"empty", "", true},
		{"whitespace only", "   ", true},
		{"too long", makeString(MaxCompanyNameLength + 1), true},
		{"valid", "Acme Corp", false},
		{"max length", makeString(MaxCompanyNameLength), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCompanyName(tt.company)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCompanyName(%q) error = %v, wantErr %v", tt.company, err, tt.wantErr)
			}
		})
	}
}

func TestValidateCompanyDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		wantErr     bool
	}{
		{"empty", "", false},
		{"too long", makeString(MaxCompanyDescriptionLength + 1), true},
		{"max length", makeString(MaxCompanyDescriptionLength), false},
		{"valid", "Company description", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCompanyDescription(tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCompanyDescription(%q) error = %v, wantErr %v", tt.description, err, tt.wantErr)
			}
		})
	}
}
