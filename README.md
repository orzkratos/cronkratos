[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/orzkratos/cronkratos/release.yml?branch=main&label=BUILD)](https://github.com/orzkratos/cronkratos/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/orzkratos/cronkratos)](https://pkg.go.dev/github.com/orzkratos/cronkratos)
[![Coverage Status](https://img.shields.io/coveralls/github/orzkratos/cronkratos/main.svg)](https://coveralls.io/github/orzkratos/cronkratos?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.25+-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/orzkratos/cronkratos.svg)](https://github.com/orzkratos/cronkratos/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/orzkratos/cronkratos)](https://goreportcard.com/report/github.com/orzkratos/cronkratos)

# cronkratos

Cron server integration with Kratos framework, wraps robfig/cron as Kratos transport.Server with clean shutdown support.

---

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->

## CHINESE README

[中文说明](README.zh.md)

<!-- TEMPLATE (EN) END: LANGUAGE NAVIGATION -->

## Main Features

🕐 **Cron Integration**: Wraps robfig/cron as Kratos transport.Server
🛡️ **Clean Shutdown**: Ensures tasks complete before exit with read-lock protection
⚡ **Context Safe**: ctx remains valid while holding read lock, prevents mid-execution invalidation
🔄 **Two Modes**: Basic mode and read-lock protected mode
📦 **Simple API**: Simple to integrate with existing Kratos applications

## Installation

```bash
go get github.com/orzkratos/cronkratos
```

## Usage

### Basic Mode (CronServer)

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
        // task logic
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

### Protected Mode (CronServerL)

Use read-lock protection to ensure ctx remains valid in execution:

```go
type MyCronServiceL struct{}

func (s *MyCronServiceL) RegisterCron(ctx context.Context, c *cron.Cron, locker sync.Locker) {
    c.AddFunc("* * * * * *", func() {
        locker.Lock()
        defer locker.Unlock()
        if ctx.Err() != nil {
            return
        }
        // task logic - ctx is guaranteed valid here
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

## Demo

See [cronkratos-demos](https://github.com/orzkratos/cronkratos-demos) as complete integration demo with Kratos applications.

## Design

The clean shutdown mechanism:

1. `cron.Stop()` - stop scheduling new tasks
2. `mutex.Lock()` - get write lock, wait read locks to release
3. `cancel()` - cancel ctx

Since cancel() is invoked following write lock acquisition, and write lock waits read locks, ctx remains valid while holding read lock. This ensures atomic operations - complete execution otherwise exit at checkpoint.

---

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-11-25 03:52:28.131064 +0000 UTC -->

## 📄 License

MIT License - see [LICENSE](LICENSE).

---

## 💬 Contact & Feedback

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Mistake reports?** Open an issue on GitHub with reproduction steps
- 💡 **Fresh ideas?** Create an issue to discuss
- 📖 **Documentation confusing?** Report it so we can improve
- 🚀 **Need new features?** Share the use cases to help us understand requirements
- ⚡ **Performance issue?** Help us optimize through reporting slow operations
- 🔧 **Configuration problem?** Ask questions about complex setups
- 📢 **Follow project progress?** Watch the repo to get new releases and features
- 🌟 **Success stories?** Share how this package improved the workflow
- 💬 **Feedback?** We welcome suggestions and comments

---

## 🔧 Development

New code contributions, follow this process:

1. **Fork**: Fork the repo on GitHub (using the webpage UI).
2. **Clone**: Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. **Navigate**: Navigate to the cloned project (`cd repo-name`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`).
5. **Code**: Implement the changes with comprehensive tests
6. **Testing**: (Golang project) Ensure tests pass (`go test ./...`) and follow Go code style conventions
7. **Documentation**: Update documentation to support client-facing changes
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) ensuring backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`).
11. **PR**: Open a merge request on GitHub (on the GitHub webpage) with detailed description.

Please ensure tests pass and include relevant documentation updates.

---

## 🌟 Support

Welcome to contribute to this project via submitting merge requests and reporting issues.

**Project Support:**

- ⭐ **Give GitHub stars** if this project helps you
- 🤝 **Share with teammates** and (golang) programming friends
- 📝 **Write tech blogs** about development tools and workflows - we provide content writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and the (golang) development scene

**Have Fun Coding with this package!** 🎉🎉🎉

<!-- TEMPLATE (EN) END: STANDARD PROJECT FOOTER -->

---

## GitHub Stars

[![Stargazers](https://starchart.cc/orzkratos/cronkratos.svg?variant=adaptive)](https://starchart.cc/orzkratos/cronkratos)
