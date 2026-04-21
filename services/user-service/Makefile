.PHONY: all build test clean docker proto

# 构建所有服务
all: build

# 构建指定服务
build-service:
	@echo "Building $(service)..."
	@cd services/$(service) && go build -o bin/server ./cmd/server

# 构建所有服务
build:
	@for service in user-service order-service payment-service gateway; do \
		echo "Building $$service..."; \
		cd services/$$service && go build -o bin/server ./cmd/server && cd ../..; \
	done

# 运行测试
test:
	@echo "Running tests..."
	@for service in user-service order-service payment-service gateway; do \
		echo "Testing $$service..."; \
		cd services/$$service && go test ./... -v && cd ../..; \
	done

# 清理
clean:
	@echo "Cleaning..."
	@for service in user-service order-service payment-service gateway; do \
		rm -rf services/$$service/bin; \
	done
	@rm -rf dist/

# Docker 构建
docker:
	@echo "Building Docker images..."
	@for service in user-service order-service payment-service gateway; do \
		docker build -t your-registry/$$service:latest -f services/$$service/deployments/docker/Dockerfile .; \
	done

# 生成 Protobuf 代码
proto:
	@echo "Generating protobuf..."
	@protoc --go_out=paths=source_relative:. \
		--go-grpc_out=paths=source_relative:. \
		--proto_path=api/proto \
		api/proto/**/*.proto

# 数据库迁移
migrate:
	@echo "Running migrations..."
	@cd services/user-service && go run scripts/migrate.go

# 代码质量检查
lint:
	@golangci-lint run ./...

# 运行所有服务
run:
	@docker-compose -f deployments/docker-compose/docker-compose.yml up