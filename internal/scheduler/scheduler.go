package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"aggregation-dashboard/internal/service"

	"github.com/robfig/cron/v3"
)

type JobStatus struct {
	SourceID  string    `json:"source_id"`
	Schedule  string    `json:"schedule"`
	LastRunAt time.Time `json:"last_run_at"`
}

type Scheduler struct {
	cron           *cron.Cron
	pollingService *service.PollingService

	mu      sync.RWMutex
	status  []JobStatus
	entries map[string]cron.EntryID
}

func NewScheduler(
	pollingService *service.PollingService,
) *Scheduler {

	return &Scheduler{
		cron: cron.New(
			cron.WithSeconds(),
		),

		pollingService: pollingService,
		status:         []JobStatus{},
		entries:        map[string]cron.EntryID{},
	}
}
func (s *Scheduler) RegisterJobs(
	cfg *SchedulerConfig,
) error {

	for _, source := range cfg.PollingSources {

		sourceCopy := source

		entryID, err := s.cron.AddFunc(
			source.Schedule,
			func() {

				log.Printf(
					"running polling job: %s",
					sourceCopy.SourceID,
				)

				err := s.pollingService.Poll(
					context.Background(),
					sourceCopy.SourceID,
					sourceCopy.URL,
					time.Now().UTC(),
				)

				if err != nil {
					log.Printf(
						"polling failed: %v",
						err,
					)
				}

				s.updateJobStatus(
					sourceCopy.SourceID,
					sourceCopy.Schedule,
				)
			},
		)

		if err != nil {
			return err
		}

		s.entries[source.SourceID] = entryID
	}

	return nil
}
func (s *Scheduler) Start() {
	s.cron.Start()
}
func (s *Scheduler) GetStatus() []JobStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.status
}
func (s *Scheduler) updateJobStatus(
	sourceID string,
	schedule string,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.status = append(
		s.status,
		JobStatus{
			SourceID:  sourceID,
			Schedule:  schedule,
			LastRunAt: time.Now(),
		},
	)
}
func (s *Scheduler) RemoveAllJobs() {
	for _, entryID := range s.entries {
		s.cron.Remove(entryID)
	}

	s.entries = map[string]cron.EntryID{}

	s.mu.Lock()
	s.status = []JobStatus{}
	s.mu.Unlock()
}
func (s *Scheduler) ReloadConfig(
	configPath string,
) error {

	log.Println("reloading scheduler config")

	cfg, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	s.RemoveAllJobs()

	err = s.RegisterJobs(cfg)
	if err != nil {
		return err
	}

	log.Println("scheduler config reloaded")

	return nil
}