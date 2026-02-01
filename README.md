# log

## 背景

本项目是 github.com/hq-devs 组织下面基于 log/slog 开发的公共日志库。其他项目中只需要引用该项目，在代码调用 `SetLogger` 即可，然后就能调用 log.Infof 方法来打印日志