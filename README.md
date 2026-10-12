# rms-lite

### 项目架构
```shell
├── bootstrap
├── cmd
│   ├── http_server
│   └── ws_server
├── core
│   ├── casbin
│   ├── config
│   ├── database
│   ├── jwt
│   ├── logger
│   └── redis
├── cron
├── docs
├── grpc
├── http
│   ├── controller
│   ├── middleware
│   ├── request
│   ├── response
│   └── route
├── model
├── repo
├── service
├── util
│   ├── crypto
│   └── formatter
└── ws
```

### 编译
```
cd /data/web/backend/rms-lite
go build -o rms_http_server cmd/http_server/main.go
go build -o rms_ws_server cmd/ws_server/main.go
```
