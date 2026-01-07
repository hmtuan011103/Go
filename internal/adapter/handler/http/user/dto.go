package user

import (
	"errors"
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
)

// Auth Requests
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *RegisterRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Email == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(r.Email) {
		return errors.New("invalid email format")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

// User Management Requests
type CreateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (r *CreateRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Email == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(r.Email) {
		return errors.New("invalid email format")
	}
	return nil
}

type UpdateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (r *UpdateRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
	if r.Name == "" && r.Email == "" {
		return errors.New("at least one field (name or email) must be provided")
	}
	if r.Email != "" && !emailRegex.MatchString(r.Email) {
		return errors.New("invalid email format")
	}
	return nil
}
