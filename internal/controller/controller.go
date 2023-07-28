package controller

// UserController interface defines the functions that are required by UserControllers
// for this project, this is only the memoryController but in the future that could be easily replaced
type UserController interface {
	GetUsers() []*User
	AddUser(*User) error
	UpdateUser(*User) error
	DeleteUser(string) error
	GetUserRoleFromHash(string) (UserRole, error)
}
