package restserver

import (
	"log"
	"net/http"

	"github.com/colbymilton/challenge-07-26-2023/internal/controller"
	"github.com/gin-gonic/gin"
)

// logMiddleware is a middleware that will log the requests being received
func logMiddleware(c *gin.Context) {
	log.Printf("Received request: %v %v\n", c.Request.Method, c.Request.RequestURI)
}

// authMiddleware is a middleware that handles authorization by ensuring that the token
// specified in the Authorization header is of a user who has one of the specified authorized roles
func authMiddleware(authorizedRoles ...controller.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get auth token from headers
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ginError("missing authorization token"))
			return
		}

		// validate the token
		// in a "real" application, this would likely be a more involved process
		// here the token is simply a hashed user email
		role, err := server.controller.GetUserRoleFromHash(authToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ginError(err.Error()))
			return
		}

		// loop through authorized roles and see if the given user's role matches.
		// in this app, there is only ever a single authorized role, so this looping
		// could have been avoided. For future-proofing, I've allowed it to handle
		// multiple authorized roles.
		authorized := false
		for _, authorizedRole := range authorizedRoles {
			if role == authorizedRole {
				authorized = true
				break
			}
		}

		if !authorized {
			c.AbortWithStatusJSON(http.StatusForbidden, ginError("insufficient permissions"))
			return
		}

		c.Next()
	}
}
