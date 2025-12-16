[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/orzkratos/cronkratos/release.yml?branch=main&label=BUILD)](https://github.com/orzkratos/cronkratos/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/orzkratos/cronkratos)](https://pkg.go.dev/github.com/orzkratos/cronkratos)
[![Coverage Status](https://img.shields.io/coveralls/github/orzkratos/cronkratos/main.svg)](https://coveralls.io/github/orzkratos/cronkratos?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.25+-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/orzkratos/cronkratos.svg)](https://github.com/orzkratos/cronkratos/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/orzkratos/cronkratos)](https://goreportcard.com/report/github.com/orzkratos/cronkratos)

# cronkratos

Cron 服务与 Kratos 框架集成，将 robfig/cron 封装成 Kratos transport.Server，支持优雅退出。

---

<!-- TEMPLATE (ZH) BEGIN: LANGUAGE NAVIGATION -->

## 英文文档

[ENGLISH README](README.md)

<!-- TEMPLATE (ZH) END: LANGUAGE NAVIGATION -->

## 核心特性

🕐 **Cron 集成**: 将 robfig/cron 封装成 Kratos transport.Server
🛡️ **优雅退出**: 通过读锁保护确保任务完成后再退出
⚡ **上下文安全**: 持有读锁时 ctx 保持有效，避免执行中途失效
🔄 **两种模式**: 基础模式和读锁保护模式
📦 **简洁 API**: 简单集成现有 Kratos 应用

## 安装

```bash
go get github.com/orzkratos/cronkratos
```

## 使用方法

### 基础模式 (CronServer)

```go
package main

import (
    "context"
    "github.com/go-kratos/kratos/v2/log"
    "github.com/orzkratos/cronkratos"
    "github.com/robfig/cron/v3"
)

type MyCronService struct{}

func (s *MyCronService) RegisterCron(ctx context.Context, c *cron.Cron) {
    c.AddFunc("* * * * * *", func() {
        // 任务逻辑
    })
}

func main() {
    c := cron.New(cron.WithSeconds())
    srv := cronkratos.NewServer(c, log.DefaultLogger)

    cronkratos.RegisterCronServer(srv, &MyCronService{})

    srv.Start(context.Background())
    // ...
    srv.Stop(context.Background())
}
```

### 保护模式 (CronServerL)

使用读锁保护确保 ctx 在执行期间保持有效：

```go
type MyCronServiceL struct{}

func (s *MyCronServiceL) RegisterCron(ctx context.Context, c *cron.Cron, locker sync.Locker) {
    c.AddFunc("* * * * * *", func() {
        locker.Lock()
        defer locker.Unlock()
        if ctx.Err() != nil {
            return
        }
        // 任务逻辑 - 此处 ctx 保证有效
    })
}

func main() {
    c := cron.New(cron.WithSeconds())
    srv := cronkratos.NewServer(c, log.DefaultLogger)

    cronkratos.RegisterCronServerL(srv, &MyCronServiceL{})

    srv.Start(context.Background())
    // ...
    srv.Stop(context.Background())
}
```

## 演示

查看 [cronkratos-demos](https://github.com/orzkratos/cronkratos-demos) 获取完整的 Kratos 应用集成演示。

## 设计说明

优雅退出机制：

1. `cron.Stop()` - 停止调度新任务
2. `mutex.Lock()` - 获取写锁，等待读锁释放
3. `cancel()` - 取消 ctx

因为 cancel() 在写锁后调用，写锁又等读锁释放，所以持有读锁时 ctx 保持有效。这确保了操作的原子性 - 要么完整执行，要么在检查点退出。

---

<!-- TEMPLATE (ZH) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-11-25 03:52:28.131064 +0000 UTC -->

## 📄 许可证类型

MIT 许可证 - 详见 [LICENSE](LICENSE)。

---

## 💬 联系与反馈

非常欢迎贡献代码！报告 BUG、建议功能、贡献代码：

- 🐛 **问题报告？** 在 GitHub 上提交问题并附上重现步骤
- 💡 **新颖思路？** 创建 issue 讨论
- 📖 **文档疑惑？** 报告问题，帮助我们完善文档
- 🚀 **需要功能？** 分享使用场景，帮助理解需求
- ⚡ **性能瓶颈？** 报告慢操作，协助解决性能问题
- 🔧 **配置困扰？** 询问复杂设置的相关问题
- 📢 **关注进展？** 关注仓库以获取新版本和功能
- 🌟 **成功案例？** 分享这个包如何改善工作流程
- 💬 **反馈意见？** 欢迎提出建议和意见

---

## 🔧 代码贡献

新代码贡献，请遵循此流程：

1. **Fork**：在 GitHub 上 Fork 仓库（使用网页界面）
2. **克隆**：克隆 Fork 的项目（`git clone https://github.com/yourname/repo-name.git`）
3. **导航**：进入克隆的项目（`cd repo-name`）
4. **分支**：创建功能分支（`git checkout -b feature/xxx`）
5. **编码**：实现您的更改并编写全面的测试
6. **测试**：（Golang 项目）确保测试通过（`go test ./...`）并遵循 Go 代码风格约定
7. **文档**：面向用户的更改需要更新文档
8. **暂存**：暂存更改（`git add .`）
9. **提交**：提交更改（`git commit -m "Add feature xxx"`）确保向后兼容的代码
10. **推送**：推送到分支（`git push origin feature/xxx`）
11. **PR**：在 GitHub 上打开 Merge Request（在 GitHub 网页上）并提供详细描述

请确保测试通过并包含相关的文档更新。

---

## 🌟 项目支持

非常欢迎通过提交 Merge Request 和报告问题来贡献此项目。

**项目支持：**

- ⭐ **给予星标**如果项目对您有帮助
- 🤝 **分享项目**给团队成员和（golang）编程朋友
- 📝 **撰写博客**关于开发工具和工作流程 - 我们提供写作支持
- 🌟 **加入生态** - 致力于支持开源和（golang）开发场景

**祝你用这个包编程愉快！** 🎉🎉🎉

<!-- TEMPLATE (ZH) END: STANDARD PROJECT FOOTER -->

---

## GitHub 标星点赞

[![标星点赞](https://starchart.cc/orzkratos/cronkratos.svg?variant=adaptive)](https://starchart.cc/orzkratos/cronkratos)
