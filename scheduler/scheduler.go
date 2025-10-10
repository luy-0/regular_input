package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/robfig/cron/v3"
)

// Task 定时任务接口
type Task interface {
	Execute(ctx context.Context) error
	GetName() string
}

// Scheduler 定时任务调度器
type Scheduler struct {
	cron    *cron.Cron
	tasks   map[string]*ScheduledTask
	mu      sync.RWMutex
	ctx     context.Context
	cancel  context.CancelFunc
	running bool
}

// ScheduledTask 已调度的任务
type ScheduledTask struct {
	Task     Task
	Schedule string
	EntryID  cron.EntryID
	Name     string
}

// NewScheduler 创建新的调度器
func NewScheduler() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		cron:    cron.New(cron.WithSeconds()), // 支持秒级定时
		tasks:   make(map[string]*ScheduledTask),
		ctx:     ctx,
		cancel:  cancel,
		running: false,
	}
}

// AddTask 添加任务到调度器
// schedule: cron 表达式，例如 "0 0 7 * * *" 表示每天早上7点执行
func (s *Scheduler) AddTask(name string, schedule string, task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查任务是否已存在
	if _, exists := s.tasks[name]; exists {
		return fmt.Errorf("任务 %s 已存在", name)
	}

	// 创建包装函数，捕获错误
	wrappedFunc := func() {
		log.Printf("[调度器] 开始执行任务: %s", name)
		if err := task.Execute(s.ctx); err != nil {
			log.Printf("[调度器] 任务执行失败 %s: %v", name, err)
		} else {
			log.Printf("[调度器] 任务执行成功: %s", name)
		}
	}

	// 添加到 cron
	entryID, err := s.cron.AddFunc(schedule, wrappedFunc)
	if err != nil {
		return fmt.Errorf("添加定时任务失败: %w", err)
	}

	// 保存任务信息
	s.tasks[name] = &ScheduledTask{
		Task:     task,
		Schedule: schedule,
		EntryID:  entryID,
		Name:     name,
	}

	log.Printf("[调度器] 已添加任务: %s, 调度表达式: %s", name, schedule)
	return nil
}

// RemoveTask 移除任务
func (s *Scheduler) RemoveTask(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[name]
	if !exists {
		return fmt.Errorf("任务 %s 不存在", name)
	}

	s.cron.Remove(task.EntryID)
	delete(s.tasks, name)
	log.Printf("[调度器] 已移除任务: %s", name)
	return nil
}

// GetTask 获取任务
func (s *Scheduler) GetTask(name string) (*ScheduledTask, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, exists := s.tasks[name]
	return task, exists
}

// GetAllTasks 获取所有任务
func (s *Scheduler) GetAllTasks() map[string]*ScheduledTask {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 返回副本避免并发问题
	tasks := make(map[string]*ScheduledTask)
	for k, v := range s.tasks {
		tasks[k] = v
	}
	return tasks
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		log.Println("[调度器] 调度器已在运行中")
		return
	}

	s.cron.Start()
	s.running = true
	log.Printf("[调度器] 调度器已启动，当前有 %d 个任务", len(s.tasks))
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		log.Println("[调度器] 调度器未在运行")
		return
	}

	// 取消所有正在执行的任务
	s.cancel()

	// 停止 cron
	cronCtx := s.cron.Stop()
	<-cronCtx.Done()

	s.running = false
	log.Println("[调度器] 调度器已停止")
}

// IsRunning 检查调度器是否在运行
func (s *Scheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// ExecuteTaskNow 立即执行指定任务（用于测试）
func (s *Scheduler) ExecuteTaskNow(name string) error {
	s.mu.RLock()
	task, exists := s.tasks[name]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("任务 %s 不存在", name)
	}

	log.Printf("[调度器] 立即执行任务: %s", name)
	if err := task.Task.Execute(s.ctx); err != nil {
		log.Printf("[调度器] 任务执行失败 %s: %v", name, err)
		return err
	}
	log.Printf("[调度器] 任务执行成功: %s", name)
	return nil
}
