package routes

import (
	"github.com/Deeksharma/taskmanager/internal/config"
	"github.com/Deeksharma/taskmanager/internal/handlers"
	"github.com/Deeksharma/taskmanager/internal/repository"
	"github.com/gin-gonic/gin"
)

// NewRouter registers all the routes
func NewRouter(taskDatabaseRepo repository.TaskDatabaseRepo) *gin.Engine {
	taskManagementHandler := handlers.NewTaskManagementHandler(taskDatabaseRepo)

	router := gin.Default()
	router.Use(gin.Recovery())
	//router.Use(middlewares.AuthMiddleware())

	g := router.Group(config.GetString("server.base_url"))

	g.GET("/tasks",
		taskManagementHandler.All)
	g.GET("/tasks/:taskId",
		taskManagementHandler.IsTaskOwner(),
		taskManagementHandler.ById)
	g.POST("/tasks",
		//taskManagementHandler.IsTaskOwner(),
		taskManagementHandler.New)
	g.PUT("/tasks/:taskId",
		//taskManagementHandler.IsTaskOwner(),
		taskManagementHandler.Update)
	g.PATCH("/tasks/:taskId",
		//taskManagementHandler.IsTaskOwner(),
		taskManagementHandler.PartialUpdate)
	g.DELETE("/tasks/:taskId",
		//taskManagementHandler.IsTaskOwner(),
		taskManagementHandler.Delete)

	return router
}
