package controller

import (
	"sync"

	"github.com/colbymilton/challenge-07-26-2023/internal/utils"
)

type memoryController struct {
	users map[string]UserRole
	mutex sync.RWMutex
}

func NewMemoryController() *memoryController {
	return &memoryController{
		users: map[string]UserRole{adminUserEmail: RoleAdmin},
	}
}

// GetUsers returns a list containing every user in the system
func (mc *memoryController) GetUsers() []*User {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	users := make([]*User, len(mc.users))
	idx := 0
	for email, role := range mc.users {
		users[idx] = &User{
			Email: email,
			Role:  role,
		}
		idx++
	}
	// sort.Slice(users, func(i, j int) bool { return users[i].Email < users[j].Email })

	return users
}

// AddUser adds the specified user to the system
// returns an error if:
// - the user email is already in use
// - if the provided user is not valid
func (mc *memoryController) AddUser(newUser *User) error {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	// check that user is valid
	if !newUser.IsValid() {
		return ErrUserNotValid
	}

	// check that user does not already exist
	if _, exists := mc.users[newUser.Email]; exists {
		return ErrUserAlreadyExists
	}

	// add user
	mc.users[newUser.Email] = newUser.Role
	return nil
}

// UpdateUser updates the specified user, using email as the unique id
// returns an error if:
// - the user does not exist
// - if the provided user is not valid
func (mc *memoryController) UpdateUser(newUser *User) error {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	// check that user is valid
	if !newUser.IsValid() {
		return ErrUserNotValid
	}

	// check that user exists
	if _, exists := mc.users[newUser.Email]; !exists {
		return ErrUserNotFound
	}

	// update user
	mc.users[newUser.Email] = newUser.Role
	return nil
}

// DeleteUser deletes the user with the specified email
// return an error if there is no user with the specified email
func (mc *memoryController) DeleteUser(email string) error {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	// check that user exists
	if _, exists := mc.users[email]; !exists {
		return ErrUserNotFound
	}

	delete(mc.users, email)
	return nil
}

// GetUserRole returns the role of the specified email hash
// returns an error if the user does not exist
func (mc *memoryController) GetUserRoleFromHash(emailHash string) (UserRole, error) {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	// this program is not optimized to find a specific user/role
	// based on a hashed email. The requirements specified that the user
	// data should be stored as email:role, so hashes are not stored anywhere.
	for email, role := range mc.users {
		if emailHash == utils.HashString(email) {
			return role, nil
		}
	}
	return "", ErrUserNotFound
}
