package enum

import "errors"

type TaskStatus string

const (
	Created    TaskStatus = "created"
	InProgress TaskStatus = "in progress"
	Succeeded  TaskStatus = "succeeded"
	Discarded  TaskStatus = "discarded"
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
		return Discarded, nil
	default:
		return "", errors.New("no such deployment status available")
	}
}

func GetAllTaskStatuses() []TaskStatus {
	return []TaskStatus{Created, InProgress, Succeeded, Discarded}
}

type RoleType string

const (
	AdminRole RoleType = "ADMIN"
	UserRole  RoleType = "USER"
)
