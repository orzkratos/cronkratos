package cronkratos

import (
	"context"
	"sync"

	"github.com/robfig/cron/v3"
)

// CronServer defines basic interface for cron job registration (no lock protection)
// 定义基础的定时任务注册接口（无锁保护）
type CronServer interface {
	RegisterCron(ctx context.Context, c *cron.Cron)
}

// RegisterCronServer registers cron jobs from CronServer to Server
// 注册 CronServer 的定时任务到 Server
func RegisterCronServer(srv *Server, svc CronServer) {
	svc.RegisterCron(srv.ctx, srv.cron)
}

// CronServerL defines interface with locker for graceful shutdown
// 定义带 locker 的接口，支持优雅退出
type CronServerL interface {
	RegisterCron(ctx context.Context, c *cron.Cron, locker sync.Locker)
}

// RegisterCronServerL registers cron jobs with read locker access
// 注册定时任务，提供读锁访问
func RegisterCronServerL(srv *Server, svc CronServerL) {
	svc.RegisterCron(srv.ctx, srv.cron, srv.mutex.RLocker())
}
