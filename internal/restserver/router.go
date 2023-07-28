package restserver

import (
	"github.com/colbymilton/challenge-07-26-2023/internal/controller"
	"github.com/gin-gonic/gin"
)

var server *UserServer

type UserServer struct {
	controller controller.UserController
}

func Run() {
	router := setupRouter()
	router.Run(":9000")
}

func setupRouter() *gin.Engine {
	if server != nil {
		panic("server is already running")
	}

	server = &UserServer{
		controller: controller.NewMemoryController(),
	}

	router := gin.Default()

	// global log middleware
	router.Use(logMiddleware)

	router.GET("/users", getUsers)

	// group endpoints that require admin authorization for auth middleware
	// the requirements state that only admin roles can POST and DELETE,
	// I've made the assumption that PATCH should also be included there.
	adminAuthGroup := router.Group("/", authMiddleware(controller.RoleAdmin))
	adminAuthGroup.POST("/users", addUser)
	adminAuthGroup.PATCH("/users", updateUser)
	adminAuthGroup.DELETE("/users/:email", deleteUser)

	return router
}
