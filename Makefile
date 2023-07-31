-include .env
export

CMD_DIR=./cmd

# run service
.PHONY: run-server
run-server:
	go run ${CMD_DIR}/server/main.go

.PHONY: build-server
build-server:
	CGO_ENABLED=0 GOARCH="amd64" GOOS=linux go build -ldflags="-s -w" -o ./bin/server ${CMD_DIR}/server/main.go

# run client
.PHONY: run-client
run-client:
	go run ${CMD_DIR}/client/main.go

.PHONY: build-client
build-client:
	CGO_ENABLED=0 GOARCH="amd64" GOOS=linux go build -ldflags="-s -w" -o ./bin/client ${CMD_DIR}/client/main.go

.PHONY: start
start:
	docker-compose up
