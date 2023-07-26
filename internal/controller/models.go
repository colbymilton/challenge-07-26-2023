package controller

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleGuest UserRole = "guest"
)

type User struct {
	Role  string `json:"role"`
	Email string `json:"email"`
}
