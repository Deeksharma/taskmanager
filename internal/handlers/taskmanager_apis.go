package handlers

import (
	"context"
	"github.com/Deeksharma/taskmanager/internal/enum"
	"github.com/Deeksharma/taskmanager/internal/errors"
	"github.com/Deeksharma/taskmanager/internal/log"
	"github.com/Deeksharma/taskmanager/internal/models"
	"github.com/Deeksharma/taskmanager/internal/repository"
	"github.com/Deeksharma/taskmanager/internal/service"
	"github.com/Deeksharma/taskmanager/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
	"strconv"
)

type TaskManagementHandler struct {
	TaskManagementService *service.TaskManagementService
}

var validate = validator.New()

func NewTaskManagementHandler(taskRepo repository.TaskDatabaseRepo) *TaskManagementHandler {
	taskManagementService := service.NewTaskManagementService(taskRepo)
	return &TaskManagementHandler{
		TaskManagementService: taskManagementService,
	}
}

// New creates a new task - in created state
func (h *TaskManagementHandler) New(c *gin.Context) {
	task := models.CreateTaskRequestBody{}
	if err := c.BindJSON(&task); err != nil {
		log.ErrorWithFields(c, map[string]interface{}{
			"error": err,
		}, "error binding request body")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTask, err := h.TaskManagementService.New(c, &task)
	if err != nil {
		log.ErrorWithFields(c, map[string]interface{}{
			"error": err,
			"task":  task,
		}, "error creating a new task")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newTask)
}

// Update updates the task
func (h *TaskManagementHandler) Update(c *gin.Context) {
	taskId := c.Param("taskId")

	task := models.UpdateTaskRequestBody{}
	if err := c.BindJSON(&task); err != nil {
		log.ErrorWithFields(c, map[string]interface{}{
			"error": err,
		}, "error binding request body")
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	validationErr := validate.Struct(task)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}
	updatedFields := make(map[string]interface{})
	updatedFields["title"] = task.Title
	updatedFields["owner"] = task.Owner
	updatedFields["description"] = task.Description
	updatedFields["status"] = task.Status

	updatedTask, err := h.TaskManagementService.Update(c, taskId, updatedFields)
	if err != nil {
		log.ErrorWithFields(c, map[string]interface{}{
			"error":   err,
			"task_id": taskId,
		}, "error updating the task")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

// PartialUpdate updates the task
func (h *TaskManagementHandler) PartialUpdate(c *gin.Context) {
	taskId := c.Param("taskId")

	task := models.Task{}
	if err := c.BindJSON(&task); err != nil {
		log.ErrorWithFields(c, map[string]interface{}{
			"error": err,
		}, "error binding request body")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedFields := make(map[string]interface{})
	if task.Title != "" {
		updatedFields["title"] = task.Title
	}
	if task.Owner != "" {
		updatedFields["owner"] = task.Owner
	}
	if task.Description != "" {
		updatedFields["description"] = task.Description
	}
	if task.Status != "" {
		updatedFields["status"] = task.Status
	}

	updatedTask, err := h.TaskManagementService.Update(c, taskId, updatedFields)
	if err != nil {
		log.ErrorWithFields(c, map[string]interface{}{
			"error":   err,
			"task_id": taskId,
		}, "error updating the task")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

// ById handler for retrieving task by id
func (h *TaskManagementHandler) ById(c *gin.Context) {
	taskId := c.Param("taskId")

	task, err := h.TaskManagementService.ById(c, taskId)

	if err != nil {
		log.ErrorWithFields(c, map[string]interface{}{
			"error":   err,
			"task_id": taskId,
		}, "error fetching task")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// All handler for retrieving full list of tasks
func (h *TaskManagementHandler) All(c *gin.Context) {
	// pagination
	pagination := make(map[string]int32)
	recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}
	page, err1 := strconv.Atoi(c.Query("page"))
	if err1 != nil || page < 1 {
		page = 1
	}
	startIndex := (page - 1) * recordPerPage
	pagination["startIndex"] = int32(startIndex)
	pagination["recordPerPage"] = int32(recordPerPage)

	// filtering
	filters := make(map[string]interface{})
	if c.Query("status") != "" {
		filters["status"] = c.Query("status")
	}
	// if user is not admin then only show their tasks
	if !h.IsAdmin(c) {
		owner, _ := c.Value("user_id").(*string)
		filters["owner"] = owner
	} else if c.Query("owner") != "" {
		filters["owner"] = c.Query("owner")
	}

	// sorting
	sort := make(map[string]interface{})
	sortBy := "updated_at"
	if c.Query("sortBy") != "" && utils.Contains(models.SortingColumns, c.Query("sortBy")) {
		sortBy = c.Query("sortBy")
	}
	sortOrder := -1
	if c.Query("order") != "" {
		if c.Query("order") == "asc" {
			sortOrder = 1
		}
	}
	sort["sortBy"] = sortBy
	sort["sortOrder"] = sortOrder
	log.InfoWithFields(c, sort, "Fetching tasks")

	tasks, err := h.TaskManagementService.All(c, filters, pagination, sort)
	if err != nil && err == mongo.ErrNoDocuments {
		log.ErrorWithFields(c, map[string]interface{}{
			"error": err,
		}, "no tasks present")
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		log.ErrorWithFields(c, map[string]interface{}{
			"error": err,
		}, "error while fetching tasks")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error() + " : error occurred while aggregating data"})
		return
	}

	c.JSON(http.StatusOK, tasks)
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": enum.Deleted,
	})
}

// IsTaskOwner middleware checks if user is owner of the task
func (h *TaskManagementHandler) IsTaskOwner() gin.HandlerFunc {
	return func(c *gin.Context) {
		taskId := c.Param("taskId")

		if !h.IsAdmin(c) && !h.IsOwner(c, taskId) {
			log.ErrorWithFields(c, map[string]interface{}{
				"task_id": taskId,
			}, "authorization failed")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": errors.ErrAuthorizationFailed})
			return
		}
		c.Next()
	}
}

// IsAdmin return true if user is Admin, else false
func (h *TaskManagementHandler) IsAdmin(c context.Context) bool {
	admin, _ := c.Value("user_type").(*string)
	return *admin == string(enum.AdminRole)
}

// IsOwner checks if the user is owner of the task
func (h *TaskManagementHandler) IsOwner(ctx context.Context, taskId string) bool {
	task, err := h.TaskManagementService.ById(ctx, taskId)
	var owner *string
	owner, _ = ctx.Value("user_id").(*string)
	log.InfoWithFields(ctx, map[string]interface{}{
		"isAdmin":  owner,
		"tasOwner": task.Owner,
	}, "task is owner")
	if err != nil {
		return false
	}
	return task.Owner == *owner
}
