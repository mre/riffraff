.DEFAULT_GOAL := build

# Build app
build:
	go build
.PHONY: build

# Clean up
clean:
	@rm -fR ./cover*
.PHONY: clean

# Run tests and generates html coverage file
cover: test
	@go tool cover -html=./coverage.text -o ./coverage.html
.PHONY: cover

# Download dependencies
depend:
	# Workaround for Go modules
	# See https://github.com/alecthomas/gometalinter/issues/521#issuecomment-415976540
	@go get -u gopkg.in/alecthomas/kingpin.v3-unstable@63abe20a23e29e80bbef8089bd3dee3ac25e5306

	@go get -u gopkg.in/alecthomas/gometalinter.v2
	@gometalinter.v2 --install
.PHONY: depend

# Install app
install:
	go install
.PHONY: install

# Run linters
lint: depend
	gometalinter.v2 \
		--disable-all \
		--exclude=vendor \
		--deadline=180s \
		--enable=gofmt \
		--linter='errch:errcheck {path}:PATH:LINE:MESSAGE' \
		--enable=errch \
		--enable=vet \
		--enable=gocyclo \
		--cyclo-over=15 \
		--enable=golint \
		--min-confidence=0.85 \
		--enable=ineffassign \
		--enable=misspell \
		./..
.PHONY: lint

# Run tests
test:
	@go test -v -race -coverprofile=./coverage.text -covermode=atomic $(shell go list ./...)
.PHONY: test
