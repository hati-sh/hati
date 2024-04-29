BUILDDIR = ./build
GOBIN = ./build/bin
GO ?= latest
GORUN = env GO111MODULE=on go run
GOBUILD = env GO111MODULE=on go build

.PHONY: build
build:
	make hati

.PHONY: dev
dev:
	@echo "\n> --- run in development mode --"
	HATI_DEBUG=true HATI_DATA_DIR=./build go run --race ./cmd/hati/main.go start --tcp --tcp-host localhost --tcp-port 4242 --rpc --rpc-host 0.0.0.0 --rpc-port 6767

.PHONY: dev-client
dev-client:
	@echo "\n> --- run in development mode --"
	HATI_DEBUG=true HATI_DATA_DIR=./build go run ./cmd/hati/main.go client --host localhost --port 4242


.PHONY: hati
hati:
	mkdir -p $(GOBIN)
	go fmt ./... && cd ./cmd/hati && $(GOBUILD) -ldflags "-w" -o ../../$(GOBIN)/hati && cd ../
	# cp config.example.yml ./build/config.yml
	# cp config.yml ./build/config.yml
	chmod +x $(GOBIN)/hati

	@echo "\n> ---"
	@echo "> Build successful. Executable in: \"$(GOBIN)/hati\" "
	@echo "> ---\n"
clean:
	rm -rf $(GOBIN)