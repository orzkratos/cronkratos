package cronkratos

import (
	"context"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"
)

// Server wraps cron instance with Kratos lifecycle
// 封装 cron 实例，支持 Kratos 生命周期管理
type Server struct {
	cron   *cron.Cron
	ctx    context.Context
	cancel context.CancelFunc
	mutex  *sync.RWMutex
	logger *log.Helper
}

// NewServer creates a new cron Server instance
// 创建新的 cron Server 实例
func NewServer(cron *cron.Cron, logger log.Logger) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		cron:   cron,
		ctx:    ctx,
		cancel: cancel,
		mutex:  &sync.RWMutex{},
		logger: log.NewHelper(logger),
	}
}

// Start implements Kratos Server interface
// 实现 Kratos Server 接口的启动方法
func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("cron server starting")
	s.cron.Start()
	return nil
}

// Stop implements Kratos Server interface
// 实现 Kratos Server 接口的停止方法
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("cron server stopping")
	// Cancel context to notify running jobs
	// 取消 context 通知运行中的任务
	s.cancel()
	// cron.Stop() returns context that is done when all running jobs complete
	// cron.Stop() 返回的 context 会在所有运行中的任务完成后取消
	stopCtx := s.cron.Stop()
	select {
	case <-stopCtx.Done():
		s.logger.Info("cron jobs stopped")
	case <-ctx.Done():
		s.logger.Warn("cron server stop timeout")
	}
	// Wait for all running jobs to complete with write lock
	// 用写锁等待所有正在执行的任务完成
	s.mutex.Lock()
	s.logger.Info("cron server stopped gracefully")
	s.mutex.Unlock()
	return nil
}
