package enum

import "errors"

type TaskStatus string

const (
	Created    TaskStatus = "Created"
	InProgress TaskStatus = "InProgress"
	Succeeded  TaskStatus = "Succeeded"
	Deleted    TaskStatus = "Deleted"
)

func GetTaskStatus(status string) (TaskStatus, error) {
	switch status {
	case "created":
		return Created, nil
	case "in progress":
		return InProgress, nil
	case "succeeded":
		return Succeeded, nil
	case "discarded":
		return Deleted, nil
	default:
		return "", errors.New("no such task status available")
	}
}

func GetAllTaskStatuses() []TaskStatus {
	return []TaskStatus{Created, InProgress, Succeeded, Deleted}
}

type RoleType string

const (
	AdminRole RoleType = "ADMIN"
	UserRole  RoleType = "USER"
)
