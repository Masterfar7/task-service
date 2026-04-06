package scheduler

import (
	"context"
	"log"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type TaskRepository interface {
	Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	GetTemplatesDueForCreation(ctx context.Context, date string) ([]taskdomain.Task, error)
}

type Scheduler struct {
	repo     TaskRepository
	interval time.Duration
	stopCh   chan struct{}
}

func New(repo TaskRepository, interval time.Duration) *Scheduler {
	return &Scheduler{
		repo:     repo,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	log.Println("Scheduler started")

	// Run immediately on start
	s.processTemplates(ctx)

	for {
		select {
		case <-ticker.C:
			s.processTemplates(ctx)
		case <-s.stopCh:
			log.Println("Scheduler stopped")
			return
		case <-ctx.Done():
			log.Println("Scheduler context cancelled")
			return
		}
	}
}

func (s *Scheduler) Stop() {
	close(s.stopCh)
}

func (s *Scheduler) processTemplates(ctx context.Context) {
	now := time.Now().UTC()
	today := now.Format("2006-01-02")

	templates, err := s.repo.GetTemplatesDueForCreation(ctx, today)
	if err != nil {
		log.Printf("Error fetching templates: %v", err)
		return
	}

	log.Printf("Processing %d templates for date %s", len(templates), today)

	for _, template := range templates {
		if err := s.createTaskFromTemplate(ctx, &template, now); err != nil {
			log.Printf("Error creating task from template %d: %v", template.ID, err)
			continue
		}

		// Update template's next occurrence
		nextOccurrence := calculateNextOccurrence(now, template.RecurrenceType, template.RecurrenceConfig)
		template.NextOccurrence = &nextOccurrence
		template.UpdatedAt = now

		if _, err := s.repo.Update(ctx, &template); err != nil {
			log.Printf("Error updating template %d: %v", template.ID, err)
		}
	}
}

func (s *Scheduler) createTaskFromTemplate(ctx context.Context, template *taskdomain.Task, now time.Time) error {
	newTask := &taskdomain.Task{
		Title:        template.Title,
		Description:  template.Description,
		Status:       taskdomain.StatusNew,
		IsTemplate:   false,
		ParentTaskID: &template.ID,
		RecurrenceType: taskdomain.RecurrenceNone,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	created, err := s.repo.Create(ctx, newTask)
	if err != nil {
		return err
	}

	log.Printf("Created task %d from template %d", created.ID, template.ID)
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
	for i := 0; i < 32; i++ { // Safety limit
		day := nextDate.Day()
		if evenOddType == "even" && day%2 == 0 {
			return nextDate
		}
		if evenOddType == "odd" && day%2 == 1 {
			return nextDate
		}
		nextDate = nextDate.AddDate(0, 0, 1)
	}
	return nextDate
}
