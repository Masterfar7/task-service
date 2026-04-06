package task

import (
	"encoding/json"
	"time"
)

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type RecurrenceType string

const (
	RecurrenceNone          RecurrenceType = "none"
	RecurrenceDaily         RecurrenceType = "daily"
	RecurrenceMonthly       RecurrenceType = "monthly"
	RecurrenceSpecificDates RecurrenceType = "specific_dates"
	RecurrenceEvenOdd       RecurrenceType = "even_odd"
)

type RecurrenceConfig struct {
	// For daily: interval in days
	Interval *int `json:"interval,omitempty"`

	// For monthly: day of month (1-30)
	DayOfMonth *int `json:"day_of_month,omitempty"`

	// For specific_dates: list of dates
	Dates []string `json:"dates,omitempty"`

	// For even_odd: "even" or "odd"
	EvenOddType *string `json:"even_odd_type,omitempty"`
}

type Task struct {
	ID               int64             `json:"id"`
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	Status           Status            `json:"status"`
	IsTemplate       bool              `json:"is_template"`
	ParentTaskID     *int64            `json:"parent_task_id,omitempty"`
	RecurrenceType   RecurrenceType    `json:"recurrence_type"`
	RecurrenceConfig *RecurrenceConfig `json:"recurrence_config,omitempty"`
	NextOccurrence   *time.Time        `json:"next_occurrence,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

func (r RecurrenceType) Valid() bool {
	switch r {
	case RecurrenceNone, RecurrenceDaily, RecurrenceMonthly, RecurrenceSpecificDates, RecurrenceEvenOdd:
		return true
	default:
		return false
	}
}

func (rc *RecurrenceConfig) ToJSON() ([]byte, error) {
	if rc == nil {
		return nil, nil
	}
	return json.Marshal(rc)
}

func RecurrenceConfigFromJSON(data []byte) (*RecurrenceConfig, error) {
	if data == nil {
		return nil, nil
	}
	var config RecurrenceConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
