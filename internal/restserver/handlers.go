package restserver

import (
	"errors"
	"net/http"

	"github.com/colbymilton/challenge-07-26-2023/internal/controller"
	"github.com/gin-gonic/gin"
)

// getUsers
func getUsers(c *gin.Context) {
	users := server.controller.GetUsers()
	c.JSON(http.StatusOK, users)
}

// addUser
func addUser(c *gin.Context) {
	// parse user from request
	var user *controller.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, ginError(err.Error()))
		return
	}

	if err := server.controller.AddUser(user); err != nil {
		handleControllerError(c, err)
		return
	} else {
		c.Status(http.StatusCreated)
	}
}

// updateUser
func updateUser(c *gin.Context) {
	// parse user from request
	var user *controller.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, ginError(err.Error()))
		return
	}

	if err := server.controller.UpdateUser(user); err != nil {
		handleControllerError(c, err)
		return
	} else {
		c.Status(http.StatusOK)
	}
}

// deleteUser
func deleteUser(c *gin.Context) {
	// get user email from request URI
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, ginError("missing email parameter"))
		return
	}

	if err := server.controller.DeleteUser(email); err != nil {
		handleControllerError(c, err)
		return
	} else {
		c.Status(http.StatusOK)
	}
}

// handleControllerError translates errors from the controller into the appropriate response code and sends the response
func handleControllerError(c *gin.Context, err error) {
	if errors.Is(err, controller.ErrUserNotFound) {
		c.JSON(http.StatusNotFound, ginError(err.Error()))

	} else if errors.Is(err, controller.ErrUserNotValid) {
		c.JSON(http.StatusBadRequest, ginError(err.Error()))

	} else if errors.Is(err, controller.ErrUserAlreadyExists) {
		c.JSON(http.StatusConflict, ginError(err.Error()))

	} else if err != nil { // catch all other errors
		c.JSON(http.StatusInternalServerError, ginError(err.Error()))
	}
}

// ginError is a simple wrapper function to turn an error string into gin.H
// used to enforce consistency with how errors are formatted ("error" key)
func ginError(err string) gin.H {
	return gin.H{"error": err}
}
