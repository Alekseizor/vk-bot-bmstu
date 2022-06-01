PWD = $(shell pwd)
NAME = vk-bot

.PHONY: run
run:
	go run $(PWD)/cmd/$(NAME)/


.PHONY: build
build:
	go build -o bin/$(NAME) $(PWD)/cmd/$(NAME)


.PHONY: test
test:
	go test $(PWD)/... -parallel=10 -coverprofile=cover.out
