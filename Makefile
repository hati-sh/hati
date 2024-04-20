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
	HATI_DEBUG=true HATI_DATA_DIR=./build go run ./cmd/hati/main.go start

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