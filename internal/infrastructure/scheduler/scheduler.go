package scheduler

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"

	"github.com/krisxia0506/bilibili-watcher/internal/application"
)

// Scheduler 管理定时任务。
type Scheduler struct {
	cronRunner  *cron.Cron
	progressSvc application.VideoProgressService
}

// NewScheduler 创建一个新的 Scheduler 实例。
func NewScheduler(progressSvc application.VideoProgressService) *Scheduler {
	// 使用支持秒字段的 cron runner
	c := cron.New(cron.WithSeconds())
	return &Scheduler{
		cronRunner:  c,
		progressSvc: progressSvc,
	}
}

// RegisterJobs 注册定时任务。
func (s *Scheduler) RegisterJobs(schedule string) error {
	// 添加获取视频进度的任务
	_, err := s.cronRunner.AddFunc(schedule, s.runRecordProgressJob)
	if err != nil {
		log.Printf("Error adding cron job with schedule '%s': %v", schedule, err)
		return err
	}
	log.Printf("Registered video progress job with schedule: %s", schedule)
	return nil
}

// runRecordProgressJob 是由 cron job 执行的函数。
func (s *Scheduler) runRecordProgressJob() {
	log.Println("Cron job starting: RecordProgressForTargetVideo")
	ctx := context.Background() // 为计划任务使用 background context
	err := s.progressSvc.RecordProgressForTargetVideo(ctx)
	if err != nil {
		// 记录错误，但任务会在下一个计划时间再次运行
		log.Printf("Error executing RecordProgressForTargetVideo job: %v", err)
	}
	log.Println("Cron job finished: RecordProgressForTargetVideo")
}

// Start 启动 cron 调度器。
func (s *Scheduler) Start() {
	log.Println("Starting cron scheduler...")
	s.cronRunner.Start()
}

// Stop 优雅地停止 cron 调度器。
func (s *Scheduler) Stop() context.Context {
	log.Println("Stopping cron scheduler...")
	return s.cronRunner.Stop()
}
