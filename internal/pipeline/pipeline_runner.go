package pipeline

import (
	"context"
	"log"
	"sync"
	"time"
)

type PipelineStatus struct {
	Status           string    `json:"status"`
	RecordsProcessed int       `json:"records_processed"`
	Errors           int       `json:"errors"`
	LastRunAt        time.Time `json:"last_run_at"`
}

type PipelineRunner struct {
	processor *Processor

	mu        sync.RWMutex
	status    PipelineStatus
	isRunning bool
}

func NewPipelineRunner(
	processor *Processor,
) *PipelineRunner {
	return &PipelineRunner{
		processor: processor,
		status: PipelineStatus{
			Status: "idle",
		},
	}
}

func (r *PipelineRunner) Run() bool {
	r.mu.Lock()

	if r.isRunning {
		r.mu.Unlock()
		return false
	}

	r.isRunning = true

	r.mu.Unlock()

	go r.runAsync()

	return true
}

func (r *PipelineRunner) runAsync() {
	r.updateStatus(func(s *PipelineStatus) {
		s.Status = "running"
		s.LastRunAt = time.Now()
		s.RecordsProcessed = 0
		s.Errors = 0
	})

	log.Println("pipeline started")

	defer func() {
		r.mu.Lock()
		r.isRunning = false
		r.mu.Unlock()
	}()

	processed, errorsCount, err := r.processor.ProcessPendingData(
		context.Background(),
	)

	if err != nil {
		log.Printf(
			"pipeline failed: %v",
			err,
		)

		r.updateStatus(func(s *PipelineStatus) {
			s.Status = "failed"
			s.Errors++
		})

		return
	}

	r.updateStatus(func(s *PipelineStatus) {
		s.Status = "completed"
		s.RecordsProcessed = processed
		s.Errors = errorsCount
	})

	log.Println("pipeline completed")
}

func (r *PipelineRunner) GetStatus() PipelineStatus {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.status
}

func (r *PipelineRunner) updateStatus(
	updateFunc func(*PipelineStatus),
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	updateFunc(&r.status)
}