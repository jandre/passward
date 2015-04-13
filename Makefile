NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
UNAME := $(shell uname -s)
ifeq ($(UNAME),Darwin)
ECHO=echo
else
ECHO=/bin/echo -e
endif

all: deps build

build:
	@mkdir -p bin/
	@$(ECHO) "$(OK_COLOR)==> Building$(NO_COLOR)"
	@go build github.com/jandre/passward/passward
	@go build github.com/jandre/passward/commands
	@go build -o bin/passward

deps:
	@$(ECHO) "$(OK_COLOR)==> Installing dependencies$(NO_COLOR)"
	@mkdir -p $(GOPATH)/src/github.com/jandre 
	@test -d  $(GOPATH)/src/github.com/jandre/passward || ln -s $(PWD) $(GOPATH)/src/github.com/jandre/passward
	@go get -d -v ./...
	@echo $(DEPS) | xargs -n1 go get -d

updatedeps:
	@$(ECHO) "$(OK_COLOR)==> Updating all dependencies$(NO_COLOR)"
	@go get -d -v -u ./...
	@echo $(DEPS) | xargs -n1 go get -d -u

clean:
	@$(ECHO) "$(OK_COLOR)==> Cleaning$(NO_COLOR)"
	@rm -rf bin/
	@find . -name ".vault" | xargs rm -rf 

format:
	go fmt ./...

test: deps
	@$(ECHO) "$(OK_COLOR)==> Testing passward...$(NO_COLOR)"
	go test -v ./...

.PHONY: all build clean deps format test updatedeps
