.PHONY: build
build:
	go build -v ./cmd/books

.PHONY: run
run: build
	./books

.PHONY: seed
seed:
	go run ./cmd/seed/main.go

.PHONY: temp
temp:
	go run ./cmd/temp/main.go

.DEFAULT_GOAL := run