# log-shipper

把分散的日志文件采集、压缩并归档为单个 `tar.gz` 包，同时构建关键词倒排索引，支持后续的离线检索（无需外网服务）。

## 用法

```bash
# 采集 example/ 下的 .log 文件，打包为 bundle.tar.gz
log-shipper collect --src example --out bundle.tar.gz

# 在归档中检索关键词 "error"
log-shipper search --bundle bundle.tar.gz --keyword error
```

目录下无 `.log` 文件、目录缺失或参数不完整时返回受控错误（退出码非 0），不会崩溃。

## 构建

```bash
go build ./...
go test ./...
```
