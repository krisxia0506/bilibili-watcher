package scheduler

import (
	"context"
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
)

// Scheduler 管理定时任务。
type Scheduler struct {
	cronRunner *cron.Cron
}

// NewScheduler 创建一个新的 Scheduler 实例。
func NewScheduler() *Scheduler {
	// 使用支持秒字段的 cron runner
	c := cron.New(cron.WithSeconds())
	return &Scheduler{
		cronRunner: c,
	}
}

// ScheduleJob 注册一个按 cronExpression 定时执行的作业。
// jobName: 用于日志记录的作业名称。
// cronExpression: cron 表达式字符串。
// job: 要执行的无参数函数。
func (s *Scheduler) ScheduleJob(jobName string, cronExpression string, job func()) error {
	entryID, err := s.cronRunner.AddFunc(cronExpression, job)
	if err != nil {
		log.Printf("Error adding cron job '%s' with schedule '%s': %v", jobName, cronExpression, err)
		return fmt.Errorf("failed to add cron job '%s': %w", jobName, err)
	}
	log.Printf("Scheduled job '%s' (EntryID: %d) with schedule: %s", jobName, entryID, cronExpression)
	return nil
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
