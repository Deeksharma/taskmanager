package handlers

import (
	"fmt"
	"github.com/Deeksharma/taskmanager/internal/enum"
	"github.com/Deeksharma/taskmanager/internal/log"
	"github.com/Deeksharma/taskmanager/internal/models"
	"github.com/Deeksharma/taskmanager/internal/repository"
	"github.com/Deeksharma/taskmanager/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
)

type TaskManagementHandler struct {
	TaskManagementService *service.TaskManagementService
}

func NewTaskManagementHandler(taskRepo repository.TaskDatabaseRepo) *TaskManagementHandler {
	taskManagementService := service.NewDeploymentService(taskRepo)
	return &TaskManagementHandler{
		TaskManagementService: taskManagementService,
	}
}

// New creates a new deployment - in deployed state - only for admins
func (h *TaskManagementHandler) New(c *gin.Context) {
	taskId := c.Param("taskId")

	task := models.Task{}
	if err := c.BindJSON(&task); err != nil {
		log.ErrorWithFields(c, map[string]interface{}{
			"error": err,
		}, "error binding request body")
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	deployment, err := h.TaskManagementService.New(c, taskId, &task)
	if err != nil {
		log.ErrorWithFields(c, map[string]interface{}{
			"error":   err,
			"task_id": taskId,
		}, "error deploying portfolio")
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, deployment)
}

// Update updates the task
func (h *TaskManagementHandler) Update(c *gin.Context) {
	taskId := c.Param("taskId")

	task := models.Task{}
	if err := c.BindJSON(&task); err != nil {
		log.ErrorWithFields(c, map[string]interface{}{
			"error": err,
		}, "error binding request body")
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	deployment, err := h.TaskManagementService.Update(c, taskId, &task)
	if err != nil {
		log.ErrorWithFields(c, map[string]interface{}{
			"error":   err,
			"task_id": taskId,
		}, "error deploying portfolio")
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, deployment)
}

// ById handler for retrieving task by id
func (h *TaskManagementHandler) ById(c *gin.Context) {
	taskId := c.Param("taskId")

	deployment, err := h.TaskManagementService.ById(c, taskId)

	if err != nil {
		log.ErrorWithFields(c, map[string]interface{}{
			"error":   err,
			"task_id": taskId,
		}, "error fetching deployment")
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, deployment)
}

// All handler for retrieving full list of tasks
func (h *TaskManagementHandler) All(c *gin.Context) {
	deployments, err := h.TaskManagementService.All(c, make(map[string]interface{}))
	if err != nil && err == mongo.ErrNoDocuments {
		log.ErrorWithFields(c, map[string]interface{}{
			"error": err,
		}, "no deployments present")
		c.AbortWithStatusJSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	} else if err != nil {
		log.ErrorWithFields(c, map[string]interface{}{
			"error": err,
		}, "error while fetching deployments")
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, deployments)
}

// Delete requests to delete the task by id
func (h *TaskManagementHandler) Delete(c *gin.Context) {
	taskId := c.Param("taskId")

	err := h.TaskManagementService.Delete(c, taskId)

	if err != nil {
		log.ErrorWithFields(c, map[string]interface{}{
			"error":   err,
			"task_id": taskId,
		}, "cannot delete task")
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": enum.Discarded,
	})
}

// IsTaskOwner middleware checks if user is owner of the task
func (h *TaskManagementHandler) IsTaskOwner() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType := c.GetString("user_type")
		userId := c.GetString("user_id")

		if userType != "ADMIN" || userId == "" { // TODO: correct the logic here
			log.ErrorWithFields(c, map[string]interface{}{
				"user_type": userType,
				"user_id":   userId,
			}, "authorization failed")
			c.AbortWithStatusJSON(http.StatusForbidden, map[string]string{"error": fmt.Errorf("authorization failed").Error()})
			return
		}
		c.Next()
	}
}
