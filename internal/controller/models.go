package controller

import "fmt"

const adminUserEmail = "admin@email.com"

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleGuest UserRole = "guest"
)

var (
	ErrUserAlreadyExists = fmt.Errorf("user already exists")
	ErrUserNotFound      = fmt.Errorf("user not found")
	ErrUserNotValid      = fmt.Errorf("user not valid")
)

func (ur UserRole) IsValid() bool {
	switch ur {
	case RoleAdmin, RoleGuest:
		return true
	default:
		return false
	}
}

type User struct {
	Email string   `json:"email"`
	Role  UserRole `json:"role"`
}

func (u *User) IsValid() bool {
	// future: could have some "actual" email validation here
	return u.Email != "" && u.Role.IsValid()
}
