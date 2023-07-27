package controller

import (
	"sync"
)

type UserController struct {
	users map[string]UserRole
	mutex sync.RWMutex
}

func NewUserController() *UserController {
	return &UserController{
		users: map[string]UserRole{adminUserEmail: RoleAdmin},
	}
}

// GetUsers returns a list containing every user in the system
func (uc *UserController) GetUsers() []*User {
	uc.mutex.RLock()
	defer uc.mutex.RUnlock()

	users := make([]*User, len(uc.users))
	idx := 0
	for email, role := range uc.users {
		users[idx] = &User{
			Email: email,
			Role:  role,
		}
		idx++
	}

	return users
}

// AddUser adds the specified user to the system
// returns an error if:
// - the user email is already in use
// - if the provided user is not valid
func (uc *UserController) AddUser(newUser *User) error {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	// check that user is valid
	if !newUser.IsValid() {
		return ErrUserNotValid
	}

	// check that user does not already exist
	if _, exists := uc.users[newUser.Email]; exists {
		return ErrUserAlreadyExists
	}

	// add user
	uc.users[newUser.Email] = newUser.Role
	return nil
}

// UpdateUser updates the specified user, using email as the unique id
// returns an error if:
// - the user does not exist
// - if the provided user is not valid
func (uc *UserController) UpdateUser(newUser *User) error {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	// check that user is valid
	if !newUser.IsValid() {
		return ErrUserNotValid
	}

	// check that user exists
	if _, exists := uc.users[newUser.Email]; !exists {
		return ErrUserNotFound
	}

	// update user
	uc.users[newUser.Email] = newUser.Role
	return nil
}

// DeleteUser deletes the user with the specified email
// return an error if there is no user with the specified email
func (uc *UserController) DeleteUser(email string) error {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	// check that user exists
	if _, exists := uc.users[email]; !exists {
		return ErrUserNotFound
	}

	delete(uc.users, email)
	return nil
}
