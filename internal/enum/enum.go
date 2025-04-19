package enum

import "errors"

type TaskStatus string

const (
	Created    TaskStatus = "Created"
	InProgress TaskStatus = "InProgress"
	Succeeded  TaskStatus = "Succeeded"
	Discarded  TaskStatus = "Discarded"
	Deleted    TaskStatus = "Deleted"
)

func GetTaskStatus(status string) (TaskStatus, error) {
	switch status {
	case "Created":
		return Created, nil
	case "InProgress":
		return InProgress, nil
	case "Succeeded":
		return Succeeded, nil
	case "Discarded":
		return Discarded, nil
	case "Deleted":
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
