package validator

import (
	"testing"
)

type TestUser struct {
	Name     string `validate:"required,min=2"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
	Bio      string `validate:"omitempty,min=10"`
}

type TestUserOptional struct {
	Name  string `validate:"omitempty,min=2"`
	Email string `validate:"omitempty,email"`
}

type TestUserNoValidation struct {
	Name  string
	Email string
}

type TestUserInvalidTag struct {
	Name string `validate:"required,min=invalid"`
}

func TestNew(t *testing.T) {
	v := New()
	if v == nil {
		t.Fatal("New() returned nil")
	}
	if v.emailRegex == nil {
		t.Fatal("emailRegex is nil")
	}
}

func TestValidator_Validate_ValidStruct(t *testing.T) {
	v := New()

	tests := []struct {
		name string
		user interface{}
	}{
		{
			name: "valid user",
			user: TestUser{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
				Bio:      "This is a valid bio with more than 10 characters",
			},
		},
		{
			name: "valid user with empty optional field",
			user: TestUser{
				Name:     "Jane Doe",
				Email:    "jane@example.com",
				Password: "password123",
				Bio:      "",
			},
		},
		{
			name: "valid user pointer",
			user: &TestUser{
				Name:     "Bob Smith",
				Email:    "bob@example.com",
				Password: "password123",
			},
		},
		{
			name: "struct with no validation tags",
			user: TestUserNoValidation{
				Name:  "Test",
				Email: "invalid-email",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.user)
			if err != nil {
				t.Errorf("Validate() error = %v, want nil", err)
			}
		})
	}
}

func TestValidator_Validate_RequiredField(t *testing.T) {
	v := New()

	tests := []struct {
		name          string
		user          interface{}
		expectedError string
	}{
		{
			name: "missing name",
			user: TestUser{
				Name:     "",
				Email:    "john@example.com",
				Password: "password123",
			},
			expectedError: "Name is required",
		},
		{
			name: "missing email",
			user: TestUser{
				Name:     "John Doe",
				Email:    "",
				Password: "password123",
			},
			expectedError: "Email is required",
		},
		{
			name: "missing password",
			user: TestUser{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "",
			},
			expectedError: "Password is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.user)
			if err == nil {
				t.Errorf("Validate() error = nil, want %v", tt.expectedError)
				return
			}
			if err.Error() != tt.expectedError {
				t.Errorf("Validate() error = %v, want %v", err.Error(), tt.expectedError)
			}
		})
	}
}

func TestValidator_Validate_EmailValidation(t *testing.T) {
	v := New()

	tests := []struct {
		name          string
		email         string
		expectedError string
	}{
		{
			name:          "invalid email - no @",
			email:         "invalid-email",
			expectedError: "Email must be a valid email",
		},
		{
			name:          "invalid email - no domain",
			email:         "user@",
			expectedError: "Email must be a valid email",
		},
		{
			name:          "invalid email - no TLD",
			email:         "user@domain",
			expectedError: "Email must be a valid email",
		},
		{
			name:          "invalid email - spaces",
			email:         "user @domain.com",
			expectedError: "Email must be a valid email",
		},
		{
			name:          "invalid email - uppercase not handled properly",
			email:         "USER@DOMAIN.COM",
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := TestUser{
				Name:     "John Doe",
				Email:    tt.email,
				Password: "password123",
			}

			err := v.Validate(user)
			if tt.expectedError == "" {
				if err != nil {
					t.Errorf("Validate() error = %v, want nil", err)
				}
			} else {
				if err == nil {
					t.Errorf("Validate() error = nil, want %v", tt.expectedError)
					return
				}
				if err.Error() != tt.expectedError {
					t.Errorf("Validate() error = %v, want %v", err.Error(), tt.expectedError)
				}
			}
		})
	}
}

func TestValidator_Validate_MinLengthValidation(t *testing.T) {
	v := New()

	tests := []struct {
		name          string
		user          TestUser
		expectedError string
	}{
		{
			name: "name too short",
			user: TestUser{
				Name:     "J",
				Email:    "john@example.com",
				Password: "password123",
			},
			expectedError: "Name must be at least 2 characters",
		},
		{
			name: "password too short",
			user: TestUser{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "pass",
			},
			expectedError: "Password must be at least 8 characters",
		},
		{
			name: "bio too short when provided",
			user: TestUser{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
				Bio:      "short",
			},
			expectedError: "Bio must be at least 10 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.user)
			if err == nil {
				t.Errorf("Validate() error = nil, want %v", tt.expectedError)
				return
			}
			if err.Error() != tt.expectedError {
				t.Errorf("Validate() error = %v, want %v", err.Error(), tt.expectedError)
			}
		})
	}
}

func TestValidator_Validate_OmitEmptyValidation(t *testing.T) {
	v := New()

	user := TestUserOptional{
		Name:  "",
		Email: "",
	}

	err := v.Validate(user)
	if err != nil {
		t.Errorf("Validate() error = %v, want nil for omitempty fields", err)
	}

	user2 := TestUserOptional{
		Name:  "J",
		Email: "invalid-email",
	}

	err = v.Validate(user2)
	if err == nil {
		t.Error("Validate() error = nil, want error for invalid field with omitempty")
	}
}

func TestValidator_Validate_NonStruct(t *testing.T) {
	v := New()

	tests := []struct {
		name  string
		input interface{}
	}{
		{
			name:  "string input",
			input: "test string",
		},
		{
			name:  "int input",
			input: 42,
		},
		{
			name:  "slice input",
			input: []string{"test"},
		},
		{
			name:  "nil input",
			input: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.input)
			if err != nil {
				t.Errorf("Validate() error = %v, want nil for non-struct input", err)
			}
		})
	}
}

func TestValidator_Validate_InvalidMinTag(t *testing.T) {
	v := New()

	user := TestUserInvalidTag{
		Name: "John Doe",
	}

	err := v.Validate(user)
	if err != nil {
		t.Errorf("Validate() error = %v, want nil for invalid min tag", err)
	}
}

func TestValidator_isEmpty(t *testing.T) {
	v := New()

	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: true,
		},
		{
			name:     "non-empty string",
			input:    "test",
			expected: false,
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: true,
		},
		{
			name:     "non-empty slice",
			input:    []string{"test"},
			expected: false,
		},
		{
			name:     "nil pointer",
			input:    (*string)(nil),
			expected: true,
		},
		{
			name:     "non-nil pointer",
			input:    &[]string{"test"}[0],
			expected: false,
		},
		{
			name:     "int (always false)",
			input:    0,
			expected: false,
		},
		{
			name:     "bool (always false)",
			input:    false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "empty string" {
				user := TestUser{Name: tt.input.(string), Email: "test@example.com", Password: "password123"}
				err := v.Validate(user)
				if tt.expected && err == nil {
					t.Error("Expected validation error for empty required field")
				}
				if !tt.expected && err != nil && err.Error() == "Name is required" {
					t.Error("Unexpected validation error for non-empty field")
				}
			}
		})
	}
}

func TestValidator_ValidEmailFormats(t *testing.T) {
	v := New()

	validEmails := []string{
		"test@example.com",
		"user.name@example.com",
		"user+tag@example.com",
		"user_name@example.com",
		"user-name@example.com",
		"123@example.com",
		"test@sub.example.com",
		"test@example.co.uk",
		"TEST@EXAMPLE.COM",
	}

	for _, email := range validEmails {
		t.Run("valid_email_"+email, func(t *testing.T) {
			user := TestUser{
				Name:     "John Doe",
				Email:    email,
				Password: "password123",
			}

			err := v.Validate(user)
			if err != nil {
				t.Errorf("Validate() error = %v, want nil for valid email %s", err, email)
			}
		})
	}
}
