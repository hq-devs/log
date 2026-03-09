# log - Go 结构化日志库

[![Go Version](https://img.shields.io/badge/go-1.25%2B-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-Apache%202.0-green)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/hq-devs/log.svg)](https://pkg.go.dev/github.com/hq-devs/log)

`log` 是一个基于 Go 标准库 `log/slog` 开发的轻量级结构化日志库，提供简洁的 API 和灵活的配置选项。

## ✨ 特性

- ✅ **基于标准库** - 使用 Go 1.21+ 内置的 `log/slog`，无第三方依赖
- ✅ **结构化日志** - 支持键值对形式的结构化日志输出
- ✅ **多日志级别** - Debug, Info, Warn, Error 四个级别
- ✅ **双 API 风格** - `Infof()` 格式化风格和 `Info()` 结构化风格
- ✅ **文件输出** - 日志自动写入文件，支持目录创建
- ✅ **源码追踪** - 自动记录调用者文件和行号
- ✅ **线程安全** - 全局日志器，并发安全
- ✅ **完整测试** - 包含单元测试，确保可靠性

## 🚀 快速开始

### 安装

```bash
go get github.com/hq-devs/log
```

### 最小示例

```go
package main

import (
    "github.com/hq-devs/log"
)

func main() {
    // 初始化日志器
    err := log.SetLogger("./logs", "app.log", "day", "info", 7, 100, 1)
    if err != nil {
        panic(err)
    }

    // 使用日志
    log.Info("应用启动", "version", "1.0.0", "pid", 12345)
    log.Infof("用户 %s 登录成功", "alice")
    log.Debug("调试信息", "key", "value")
    log.Warn("磁盘空间不足", "usage", "95%")
    log.Error("数据库连接失败", "err", err)
}
```

## 📖 完整示例

### 基础配置

```go
package main

import (
    "os"
    "github.com/hq-devs/log"
)

func main() {
    // 配置日志
    err := log.SetLogger(
        "./logs",           // 日志目录
        "app.log",          // 日志文件名
        "day",              // 轮转单位 (预留参数)
        "debug",            // 日志级别: debug, info, warn, error
        7,                  // 保留天数 (预留参数)
        100,                // 单个文件大小MB (预留参数)
        1,                  // 压缩阈值 (预留参数)
    )
    if err != nil {
        log.Error("日志初始化失败", "err", err)
        os.Exit(1)
    }

    // 业务代码
    processOrder("order-123")
}

func processOrder(orderID string) {
    log.Info("开始处理订单", "order_id", orderID, "step", "validation")
    
    // 模拟业务逻辑
    if err := validateOrder(orderID); err != nil {
        log.Error("订单验证失败", "order_id", orderID, "err", err)
        return
    }
    
    log.Info("订单处理完成", "order_id", orderID, "status", "success")
}

func validateOrder(orderID string) error {
    log.Debug("验证订单", "order_id", orderID, "method", "validateOrder")
    // 验证逻辑...
    return nil
}
```

### 日志输出示例

```
time=2026-03-09T17:30:00.000Z level=INFO msg="应用启动" version=1.0.0 pid=12345 source=main/main.go:15
time=2026-03-09T17:30:01.000Z level=INFO msg="用户 alice 登录成功" source=main/main.go:16
time=2026-03-09T17:30:02.000Z level=DEBUG msg="调试信息" key=value source=utils/helper.go:42
time=2026-03-09T17:30:03.000Z level=WARN msg="磁盘空间不足" usage=95% source=monitor/disk.go:78
```

## 🔧 API 参考

### 初始化函数

```go
// SetLogger 初始化全局日志器
// 参数:
//   dir: 日志文件目录
//   file: 日志文件名
//   unit: 日志轮转单位 (预留，当前未实现)
//   level: 日志级别 "debug", "info", "warn", "error"
//   count: 保留文件数量 (预留，当前未实现)
//   size: 单个文件大小MB (预留，当前未实现)
//   compressT: 压缩阈值 (预留，当前未实现)
func SetLogger(dir, file, unit, level string, count, size, compressT int64) error
```

### 日志函数

#### 格式化风格 (Printf 风格)
```go
log.Debugf(format string, a ...interface{})
log.Infof(format string, a ...interface{})
log.Warnf(format string, a ...interface{})
log.Errorf(format string, a ...interface{})
```

#### 结构化风格 (键值对)
```go
log.Debug(msg string, v ...interface{})
log.Info(msg string, v ...interface{})
log.Warn(msg string, v ...interface{})
log.Error(msg string, v ...interface{})
```

### 使用示例对比

```go
// 格式化风格 - 适合简单消息
log.Infof("用户 %s 在 %s 登录", username, time.Now().Format("2006-01-02"))

// 结构化风格 - 适合复杂数据
log.Info("用户登录",
    "username", username,
    "login_time", time.Now(),
    "ip", "192.168.1.100",
    "user_agent", req.UserAgent(),
)
```

## ⚙️ 配置说明

### 日志级别

| 级别 | 描述 | 输出内容 |
|------|------|----------|
| debug | 调试信息 | debug, info, warn, error |
| info | 一般信息 | info, warn, error |
| warn | 警告信息 | warn, error |
| error | 错误信息 | error |

### 源码追踪

日志会自动记录调用者的文件和行号，格式为 `目录名/文件名:行号`：

```
source=log/log_test.go:25      // 测试文件
source=api/user.go:42          // 业务代码
source=utils/validator.go:18   // 工具函数
```

## 🧪 测试

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行测试并显示详细信息
go test -v ./...

# 运行特定测试
go test -v -run TestLog
```

### 测试示例

查看 [log_test.go](log_test.go) 了解完整的测试用例。

## 📁 项目结构

```
github.com/hq-devs/log/
├── .github/          # GitHub 工作流配置
├── .gitignore        # Git 忽略文件
├── LICENSE          # Apache 2.0 许可证
├── README.md        # 本文档
├── go.mod          # Go 模块定义
├── log.go          # 主实现文件
└── log_test.go     # 单元测试文件
```

## 🔄 开发计划

### 已实现
- [x] 基础日志功能
- [x] 多日志级别支持
- [x] 结构化日志输出
- [x] 文件路径压缩
- [x] 单元测试

### 计划中
- [ ] 日志轮转功能
- [ ] JSON 格式输出
- [ ] 控制台输出支持
- [ ] 上下文支持
- [ ] 性能优化
- [ ] 更完善的文档

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启 Pull Request

## 📄 许可证

本项目基于 Apache License 2.0 许可证开源，详见 [LICENSE](LICENSE) 文件。

## 📞 支持

- 提交 Issue: [GitHub Issues](https://github.com/hq-devs/log/issues)
- 组织主页: [hq-devs](https://github.com/hq-devs)

---

**提示**: 当前版本为基础版本，部分配置参数为预留功能，将在后续版本中实现。
```

## 使用建议

1. **生产环境**: 建议使用 `info` 或 `warn` 级别
2. **开发环境**: 可以使用 `debug` 级别查看更多信息
3. **日志目录**: 确保应用有写入权限
4. **性能考虑**: 高频日志调用可能影响性能，建议合理使用日志级别

## 常见问题

### Q: 为什么日志文件没有创建？
A: 检查目录权限和路径是否正确，确保应用有写入权限。

### Q: 如何更改日志格式？
A: 当前版本仅支持 Text 格式，JSON 格式支持正在开发中。

### Q: 日志级别不生效？
A: 确保在 `SetLogger` 中正确设置级别参数，注意大小写不敏感。

### Q: 如何禁用日志？
A: 设置日志级别为高于实际使用的级别，或不在代码中调用 `SetLogger`。

---

*由 [hq-devs](https://github.com/hq-devs) 组织维护 | 最后更新: 2026-03-09*