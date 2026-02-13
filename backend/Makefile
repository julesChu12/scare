GO_AIR := $(shell go env GOPATH)/bin/air

.PHONY: dev serve build test clean

# 热重载开发
dev:
	$(GO_AIR)

# 标准启动
serve:
	go run . serve

# 编译
build:
	go build -o ./tmp/main .

# 测试
test:
	go test ./...

# 清理
clean:
	rm -rf ./tmp
