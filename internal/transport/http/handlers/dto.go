package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type taskMutationDTO struct {
	Title            string                       `json:"title"`
	Description      string                       `json:"description"`
	Status           taskdomain.Status            `json:"status"`
	IsTemplate       *bool                        `json:"is_template,omitempty"`
	RecurrenceType   *taskdomain.RecurrenceType   `json:"recurrence_type,omitempty"`
	RecurrenceConfig *taskdomain.RecurrenceConfig `json:"recurrence_config,omitempty"`
}

type taskDTO struct {
	ID               int64                        `json:"id"`
	Title            string                       `json:"title"`
	Description      string                       `json:"description"`
	Status           taskdomain.Status            `json:"status"`
	IsTemplate       bool                         `json:"is_template"`
	ParentTaskID     *int64                       `json:"parent_task_id,omitempty"`
	RecurrenceType   taskdomain.RecurrenceType    `json:"recurrence_type"`
	RecurrenceConfig *taskdomain.RecurrenceConfig `json:"recurrence_config,omitempty"`
	NextOccurrence   *time.Time                   `json:"next_occurrence,omitempty"`
	CreatedAt        time.Time                    `json:"created_at"`
	UpdatedAt        time.Time                    `json:"updated_at"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:               task.ID,
		Title:            task.Title,
		Description:      task.Description,
		Status:           task.Status,
		IsTemplate:       task.IsTemplate,
		ParentTaskID:     task.ParentTaskID,
		RecurrenceType:   task.RecurrenceType,
		RecurrenceConfig: task.RecurrenceConfig,
		NextOccurrence:   task.NextOccurrence,
		CreatedAt:        task.CreatedAt,
		UpdatedAt:        task.UpdatedAt,
	}
}
