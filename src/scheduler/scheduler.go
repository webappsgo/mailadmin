package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"
	
	"github.com/webappsgo/mail/src/database"
	"github.com/webappsgo/mail/src/logger"
)

// Scheduler manages background tasks per AI.md PART 19
type Scheduler struct {
	db     *database.DB
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	
	tasks map[string]*Task
	mu    sync.RWMutex
}

// Task represents a scheduled task
type Task struct {
	ID          string
	Name        string
	Schedule    string
	Handler     TaskHandler
	Enabled     bool
	LastRun     time.Time
	NextRun     time.Time
	RunCount    int64
	FailCount   int64
	LastStatus  string
	LastError   string
}

// TaskHandler is the function signature for task execution
type TaskHandler func(ctx context.Context) error

// New creates a new scheduler instance per AI.md PART 19
func New(db *database.DB) (*Scheduler, error) {
	ctx, cancel := context.WithCancel(context.Background())
	
	s := &Scheduler{
		db:     db,
		ctx:    ctx,
		cancel: cancel,
		tasks:  make(map[string]*Task),
	}
	
	// Load task state from database
	if err := s.loadState(); err != nil {
		cancel()
		return nil, fmt.Errorf("loading scheduler state: %w", err)
	}
	
	return s, nil
}

// Start starts the scheduler per AI.md PART 19 Step 16
func (s *Scheduler) Start() {
	logger.Info("Scheduler starting...")
	
	// Register built-in tasks per AI.md PART 19
	s.registerBuiltInTasks()
	
	// Check for missed tasks (catch-up)
	s.checkMissedTasks()
	
	// Start scheduler loop
	s.wg.Add(1)
	go s.run()
	
	logger.Info("Scheduler started successfully")
}

// Stop stops the scheduler gracefully per AI.md PART 19
func (s *Scheduler) Stop() {
	logger.Info("Scheduler stopping...")
	
	// Signal shutdown
	s.cancel()
	
	// Wait for running tasks (max 30 seconds per AI.md PART 19)
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		logger.Info("All scheduled tasks completed")
	case <-time.After(30 * time.Second):
		logger.Warn("Scheduler shutdown timeout - forcing stop")
	}
	
	// Save final state
	if err := s.saveState(); err != nil {
		logger.Errorf("Failed to save scheduler state: %v", err)
	}
	
	logger.Info("Scheduler stopped")
}

// run is the main scheduler loop
func (s *Scheduler) run() {
	defer s.wg.Done()
	
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.checkAndRunTasks()
		}
	}
}

// checkAndRunTasks checks all tasks and runs those that are due
func (s *Scheduler) checkAndRunTasks() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	now := time.Now()
	
	for id, task := range s.tasks {
		if !task.Enabled {
			continue
		}
		
		// Check if task is due
		if now.After(task.NextRun) || now.Equal(task.NextRun) {
			// Run task in goroutine
			s.wg.Add(1)
			go func(t *Task) {
				defer s.wg.Done()
				s.runTask(t)
			}(task)
			
			// Update next run time
			s.updateNextRun(id)
		}
	}
}

// runTask executes a single task
func (s *Scheduler) runTask(task *Task) {
	logger.Infof("Running scheduled task: %s", task.Name)
	
	startTime := time.Now()
	
	// Execute task handler
	err := task.Handler(s.ctx)
	
	duration := time.Since(startTime)
	
	// Update task status
	s.mu.Lock()
	task.LastRun = startTime
	if err != nil {
		task.LastStatus = "failed"
		task.LastError = err.Error()
		task.FailCount++
		logger.Errorf("Task %s failed: %v (duration: %v)", task.Name, err, duration)
	} else {
		task.LastStatus = "success"
		task.LastError = ""
		task.RunCount++
		logger.Infof("Task %s completed successfully (duration: %v)", task.Name, duration)
	}
	s.mu.Unlock()
	
	// Save task state to database
	s.saveTaskState(task)
	
	// Audit log per AI.md PART 11
	logger.Audit("task.executed", "system", map[string]interface{}{
		"task":     task.Name,
		"status":   task.LastStatus,
		"duration": duration.String(),
	})
}

// registerBuiltInTasks registers all required tasks per AI.md PART 19
func (s *Scheduler) registerBuiltInTasks() {
	// Session cleanup - every 15 minutes per AI.md PART 19
	s.RegisterTask("session_cleanup", "Session Cleanup", "@every 15m", s.taskSessionCleanup, true)
	
	// Token cleanup - every 15 minutes per AI.md PART 19
	s.RegisterTask("token_cleanup", "Token Cleanup", "@every 15m", s.taskTokenCleanup, true)
	
	// Health check - every 5 minutes per AI.md PART 19
	s.RegisterTask("healthcheck_self", "Self Health Check", "@every 5m", s.taskHealthCheck, true)
	
	// SSL renewal - daily at 03:00 per AI.md PART 19
	s.RegisterTask("ssl_renewal", "SSL Certificate Renewal", "0 3 * * *", s.taskSSLRenewal, true)
	
	// GeoIP update - weekly Sunday at 03:00 per AI.md PART 19
	s.RegisterTask("geoip_update", "GeoIP Database Update", "0 3 * * 0", s.taskGeoIPUpdate, true)
	
	// Log rotation - daily at 00:00 per AI.md PART 19
	s.RegisterTask("log_rotation", "Log Rotation", "0 0 * * *", s.taskLogRotation, true)
	
	// Daily backup - daily at 02:00 per AI.md PART 19
	s.RegisterTask("backup_daily", "Daily Backup", "0 2 * * *", s.taskBackupDaily, true)
}

// RegisterTask registers a new task
func (s *Scheduler) RegisterTask(id, name, schedule string, handler TaskHandler, enabled bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	task := &Task{
		ID:       id,
		Name:     name,
		Schedule: schedule,
		Handler:  handler,
		Enabled:  enabled,
		NextRun:  s.calculateNextRun(schedule),
	}
	
	s.tasks[id] = task
	logger.Infof("Registered task: %s (schedule: %s, enabled: %v)", name, schedule, enabled)
}

// calculateNextRun calculates the next run time based on schedule
// Per AI.md PART 19: Supports cron and interval formats
func (s *Scheduler) calculateNextRun(schedule string) time.Time {
	now := time.Now()
	
	// Parse interval format (@every Xm, @every Xh)
	if len(schedule) > 7 && schedule[:7] == "@every " {
		duration, err := time.ParseDuration(schedule[7:])
		if err == nil {
			return now.Add(duration)
		}
	}
	
	// Parse special formats
	switch schedule {
	case "@hourly":
		return now.Add(1 * time.Hour).Truncate(time.Hour)
	case "@daily":
		return now.Add(24 * time.Hour).Truncate(24 * time.Hour)
	case "@weekly":
		// Next Sunday at 00:00
		days := (7 - int(now.Weekday())) % 7
		if days == 0 {
			days = 7
		}
		return now.Add(time.Duration(days) * 24 * time.Hour).Truncate(24 * time.Hour)
	case "@monthly":
		// First day of next month at 00:00
		year, month, _ := now.Date()
		return time.Date(year, month+1, 1, 0, 0, 0, 0, now.Location())
	}
	
	// TODO: Parse full cron format (0 2 * * *)
	// For now, default to 1 hour from now
	return now.Add(1 * time.Hour)
}

// updateNextRun updates the next run time for a task
func (s *Scheduler) updateNextRun(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if task, ok := s.tasks[id]; ok {
		task.NextRun = s.calculateNextRun(task.Schedule)
	}
}

// checkMissedTasks checks for and runs missed tasks per AI.md PART 19
func (s *Scheduler) checkMissedTasks() {
	// TODO: Implement catch-up logic with catch_up_window (1h default)
	logger.Info("Checking for missed tasks...")
}

// loadState loads scheduler state from database
func (s *Scheduler) loadState() error {
	// TODO: Load task state from scheduler_tasks table per AI.md PART 10
	return nil
}

// saveState saves scheduler state to database
func (s *Scheduler) saveState() error {
	// TODO: Save task state to scheduler_tasks table per AI.md PART 10
	return nil
}

// saveTaskState saves a single task's state to database
func (s *Scheduler) saveTaskState(task *Task) {
	// TODO: Update scheduler_history table per AI.md PART 10
}

// Built-in task handlers per AI.md PART 19

func (s *Scheduler) taskSessionCleanup(ctx context.Context) error {
	// Remove expired admin sessions per AI.md PART 19
	result, err := s.db.ServerDB.ExecContext(ctx, `
		DELETE FROM admin_sessions 
		WHERE expires_at < ?
	`, time.Now().Unix())
	
	if err != nil {
		return fmt.Errorf("cleaning admin sessions: %w", err)
	}
	
	rows, _ := result.RowsAffected()
	if rows > 0 {
		logger.Infof("Cleaned up %d expired admin sessions", rows)
	}
	
	return nil
}

func (s *Scheduler) taskTokenCleanup(ctx context.Context) error {
	// Remove expired tokens per AI.md PART 19
	// TODO: Implement when tokens table exists
	return nil
}

func (s *Scheduler) taskHealthCheck(ctx context.Context) error {
	// Self health verification per AI.md PART 19
	// TODO: Implement health checks
	return nil
}

func (s *Scheduler) taskSSLRenewal(ctx context.Context) error {
	// Renew SSL certificates 7 days before expiry per AI.md PART 19
	// TODO: Implement when SSL support is added (PART 15)
	return nil
}

func (s *Scheduler) taskGeoIPUpdate(ctx context.Context) error {
	// Update GeoIP databases per AI.md PART 19
	// TODO: Implement when GeoIP support is added (PART 20)
	return nil
}

func (s *Scheduler) taskLogRotation(ctx context.Context) error {
	// Rotate and compress old logs per AI.md PART 19
	// TODO: Implement log rotation
	return nil
}

func (s *Scheduler) taskBackupDaily(ctx context.Context) error {
	// Full backup + daily incremental per AI.md PART 19
	// TODO: Implement when backup system is added (PART 22)
	return nil
}
