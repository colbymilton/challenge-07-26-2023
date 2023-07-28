package controller

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func initController(userCount int) *memoryController {
	controller := NewMemoryController()
	for i := 1; i < userCount; i++ {
		controller.users[fmt.Sprintf("%d@email.com", i)] = RoleGuest
	}
	return controller
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestGetUsers(t *testing.T) {
	// get a controller with 5 existing users
	controller := initController(5)

	// confirm all users are returned
	require.Equal(t, len(controller.users), len(controller.GetUsers()))
}

func TestAddUser(t *testing.T) {
	// get a controller with 1 existing user
	controller := initController(1)

	// add a new user
	user := &User{Email: "test@email.com", Role: RoleGuest}
	require.NoError(t, controller.AddUser(user))

	// attempt to add the same user
	require.ErrorIs(t, controller.AddUser(user), ErrUserAlreadyExists)

	// attempt to add invalid user
	user.Role = "invalid"
	require.ErrorIs(t, controller.AddUser(user), ErrUserNotValid)

	// confirm user count
	require.Len(t, controller.users, 2)
}

func TestUpdateUser(t *testing.T) {
	// get a controller with 1 existing user
	controller := initController(1)

	// update the user
	user := &User{Email: adminUserEmail, Role: RoleGuest}
	require.NoError(t, controller.UpdateUser(user))

	// attempt to update missing user
	user.Email = "fake@email.com"
	require.ErrorIs(t, controller.UpdateUser(user), ErrUserNotFound)

	// attempt to update to invalid user
	user.Role = "invalid"
	require.ErrorIs(t, controller.UpdateUser(user), ErrUserNotValid)

	// confirm that user is updated
	require.Equal(t, RoleGuest, controller.users[adminUserEmail])
}

func TestDeleteUser(t *testing.T) {
	// get a controller with 5 existing users
	controller := initController(5)

	// delete a user
	require.NoError(t, controller.DeleteUser(adminUserEmail))

	// attempt to delete missing user
	require.ErrorIs(t, controller.DeleteUser("fake@email.com"), ErrUserNotFound)

	// confirm that a user was deleted
	require.Len(t, controller.users, 4)
}

func TestConcurrency(t *testing.T) {
	controller := initController(1)

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(2)
		go concurrentAdd(controller, fmt.Sprintf("%v@email.com", i), &wg)
		go concurrentRead(controller, &wg)
	}
	wg.Wait()
	require.Len(t, controller.users, 1001)
}

func concurrentAdd(controller *memoryController, email string, wg *sync.WaitGroup) {
	controller.AddUser(&User{Email: email, Role: RoleGuest})
	wg.Done()
}

func concurrentRead(controller *memoryController, wg *sync.WaitGroup) {
	controller.GetUsers()
	wg.Done()
}
