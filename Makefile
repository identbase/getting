# Commands
GO_CMD = `which go`
LINT_CMD = $(GOPATH)/bin/golint

# Directories
PACKAGE = github.com/identbase/getting
SRC = $(GOPATH)/src/$(PACKAGE)

default: lint test

lint:
	$(LINT_CMD) ./...

test:
	cd $(SRC)
	$(GO_CMD) test -v ./...
	# $(GOPATH)/bin/goveralls -coverprofile=coverage.out -service=travis-ci # -repotoken $(COVERALLS_TOKEN)

# coverage:
# 	cd $(SRC)
# 	$(GOPATH)/bin/overalls -project=$(PACKAGE) -covermode=count
# 	$(GOPATH)/bin/goveralls -coverprofile=overalls.coverprofile -service=travis-ci
