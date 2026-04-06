package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	now := s.now()
	model := &taskdomain.Task{
		Title:            normalized.Title,
		Description:      normalized.Description,
		Status:           normalized.Status,
		IsTemplate:       normalized.IsTemplate,
		RecurrenceType:   normalized.RecurrenceType,
		RecurrenceConfig: normalized.RecurrenceConfig,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	// Calculate next occurrence for templates
	if model.IsTemplate && model.RecurrenceType != taskdomain.RecurrenceNone {
		nextOccurrence := calculateNextOccurrence(now, model.RecurrenceType, model.RecurrenceConfig)
		model.NextOccurrence = &nextOccurrence
	}

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	now := s.now()
	model := &taskdomain.Task{
		ID:               id,
		Title:            normalized.Title,
		Description:      normalized.Description,
		Status:           normalized.Status,
		IsTemplate:       normalized.IsTemplate,
		RecurrenceType:   normalized.RecurrenceType,
		RecurrenceConfig: normalized.RecurrenceConfig,
		UpdatedAt:        now,
	}

	// Recalculate next occurrence for templates
	if model.IsTemplate && model.RecurrenceType != taskdomain.RecurrenceNone {
		nextOccurrence := calculateNextOccurrence(now, model.RecurrenceType, model.RecurrenceConfig)
		model.NextOccurrence = &nextOccurrence
	}

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if input.RecurrenceType == "" {
		input.RecurrenceType = taskdomain.RecurrenceNone
	}

	if !input.RecurrenceType.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid recurrence type", ErrInvalidInput)
	}

	// Validate recurrence config
	if input.RecurrenceType != taskdomain.RecurrenceNone {
		if err := validateRecurrenceConfig(input.RecurrenceType, input.RecurrenceConfig); err != nil {
			return CreateInput{}, err
		}
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if input.RecurrenceType == "" {
		input.RecurrenceType = taskdomain.RecurrenceNone
	}

	if !input.RecurrenceType.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid recurrence type", ErrInvalidInput)
	}

	// Validate recurrence config
	if input.RecurrenceType != taskdomain.RecurrenceNone {
		if err := validateRecurrenceConfig(input.RecurrenceType, input.RecurrenceConfig); err != nil {
			return UpdateInput{}, err
		}
	}

	return input, nil
}

func validateRecurrenceConfig(recType taskdomain.RecurrenceType, config *taskdomain.RecurrenceConfig) error {
	if config == nil {
		return fmt.Errorf("%w: recurrence_config is required for recurrence type %s", ErrInvalidInput, recType)
	}

	switch recType {
	case taskdomain.RecurrenceDaily:
		if config.Interval == nil || *config.Interval < 1 {
			return fmt.Errorf("%w: interval must be >= 1 for daily recurrence", ErrInvalidInput)
		}
	case taskdomain.RecurrenceMonthly:
		if config.DayOfMonth == nil || *config.DayOfMonth < 1 || *config.DayOfMonth > 30 {
			return fmt.Errorf("%w: day_of_month must be between 1 and 30", ErrInvalidInput)
		}
	case taskdomain.RecurrenceSpecificDates:
		if len(config.Dates) == 0 {
			return fmt.Errorf("%w: dates array cannot be empty for specific_dates recurrence", ErrInvalidInput)
		}
	case taskdomain.RecurrenceEvenOdd:
		if config.EvenOddType == nil || (*config.EvenOddType != "even" && *config.EvenOddType != "odd") {
			return fmt.Errorf("%w: even_odd_type must be 'even' or 'odd'", ErrInvalidInput)
		}
	}

	return nil
}

func calculateNextOccurrence(from time.Time, recType taskdomain.RecurrenceType, config *taskdomain.RecurrenceConfig) time.Time {
	switch recType {
	case taskdomain.RecurrenceDaily:
		return from.AddDate(0, 0, *config.Interval)
	case taskdomain.RecurrenceMonthly:
		return getNextMonthlyDate(from, *config.DayOfMonth)
	case taskdomain.RecurrenceSpecificDates:
		return getNextSpecificDate(from, config.Dates)
	case taskdomain.RecurrenceEvenOdd:
		return getNextEvenOddDate(from, *config.EvenOddType)
	default:
		return from
	}
}

func getNextMonthlyDate(from time.Time, dayOfMonth int) time.Time {
	year, month, _ := from.Date()
	nextDate := time.Date(year, month, dayOfMonth, 0, 0, 0, 0, time.UTC)

	if nextDate.Before(from) || nextDate.Equal(from) {
		nextDate = nextDate.AddDate(0, 1, 0)
	}

	// Handle months with fewer days
	if nextDate.Day() != dayOfMonth {
		nextDate = nextDate.AddDate(0, 1, 0)
		nextDate = time.Date(nextDate.Year(), nextDate.Month(), dayOfMonth, 0, 0, 0, 0, time.UTC)
	}

	return nextDate
}

func getNextSpecificDate(from time.Time, dates []string) time.Time {
	var nextDate time.Time
	for _, dateStr := range dates {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		if parsed.After(from) && (nextDate.IsZero() || parsed.Before(nextDate)) {
			nextDate = parsed
		}
	}
	if nextDate.IsZero() {
		return from.AddDate(100, 0, 0) // Far future if no valid dates
	}
	return nextDate
}

func getNextEvenOddDate(from time.Time, evenOddType string) time.Time {
	nextDate := from.AddDate(0, 0, 1)
	for {
		day := nextDate.Day()
		if evenOddType == "even" && day%2 == 0 {
			return nextDate
		}
		if evenOddType == "odd" && day%2 == 1 {
			return nextDate
		}
		nextDate = nextDate.AddDate(0, 0, 1)
		// Safety: don't loop forever
		if nextDate.After(from.AddDate(0, 1, 0)) {
			break
		}
	}
	return nextDate
}
