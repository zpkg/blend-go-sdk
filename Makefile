PREFIX			?= $(shell pwd)
PKGS 			:= $(shell go list ./... | grep -v /vendor/)
SHASUMCMD 		:= $(shell command -v sha1sum || command -v shasum; 2> /dev/null)
TARCMD 			:= $(shell command -v tar || command -v tar; 2> /dev/null)
GIT_REF 		:= $(shell git log --pretty=format:'%h' -n 1)
CURRENT_USER 	:= $(shell whoami)
VERSION 		:= $(shell cat ./VERSION)

# this is to allow local go-sdk/db tests to pass
DB_PORT 		?= 5432
DB_SSLMODE		?= disable

# coverage stuff
CIRCLE_ARTIFACTS 	?= "."
COVERAGE_OUT 		:= "$(CIRCLE_ARTIFACTS)/coverage.html"

export GIT_REF
export VERSION
export DB_SSLMODE

all: format vet profanity test

ci: vet profanity cover 

new-install:
	@go get -u github.com/lib/pq
	@go get -u github.com/airbrake/gobrake
	@go get -u github.com/DataDog/datadog-go/statsd
	@go get -u github.com/opentracing/opentracing-go
	@go get -u golang.org/x/net/http2
	@go get -u golang.org/x/oauth2
	@go get -u golang.org/x/oauth2/google
	@go get -u golang.org/x/lint/golint

format:
	@echo "$(VERSION)/$(GIT_REF) >> formatting code"
	@go fmt $(PKGS)

vet:
	@echo "$(VERSION)/$(GIT_REF) >> vetting code"
	@go vet $(PKGS)

lint:
	@echo "$(VERSION)/$(GIT_REF) >> linting code"
	@golint $(PKGS)

build:
	@echo "$(VERSION)/$(GIT_REF) >> linting code"
	@docker build . -t go-sdk:$(GIT_REF) 
	@docker run go-sdk:$(GIT_REF)

.PHONY: profanity
profanity:
	@echo "$(VERSION)/$(GIT_REF) >> profanity"
	@go run _bin/profanity/main.go -rules PROFANITY --exclude="_bin/*"

test-circleci:
	@echo "$(VERSION)/$(GIT_REF) >> tests"
	@circleci build

test:
	@echo "$(VERSION)/$(GIT_REF) >> tests"
	@go test $(PKGS) -timeout 15s

test-verbose:
	@echo "$(VERSION)/$(GIT_REF) >> tests"
	@go test -v $(PKGS)

cover:
	@echo "$(VERSION)/$(GIT_REF) >> coverage"
	@go run _bin/coverage/main.go

cover-enforce:
	@echo "$(VERSION)/$(GIT_REF) >> coverage"
	@go run _bin/coverage/main.go -enforce
	
cover-update:
	@echo "$(VERSION)/$(GIT_REF) >> coverage"
	@go run _bin/coverage/main.go -update

increment-patch:
	@echo "Current Version $(VERSION)"
	@go run _bin/semver/main.go patch ./.version > ./NEW_VERSION
	@mv ./NEW_VERSION ./.version
	@cat ./.version

increment-minor:
	@echo "Current Version $(VERSION)"
	@go run _bin/semver/main.go minor ./.version > ./NEW_VERSION
	@mv ./NEW_VERSION ./.version
	@cat ./.version

increment-major:
	@echo "Current Version $(VERSION)"
	@go run _bin/semver/main.go major ./VERSION > ./NEW_VERSION
	@mv ./NEW_VERSION ./VERSION
	@cat ./VERSION

clean:
	@echo "Cleaning COVERAGE files"
	@find . -name "COVERAGE" -exec rm {} \;

tag:
	git tag -f $(VERSION)

push-tags:
	git push -f --tags

install-profanity:
	@go install github.com/blend/go-sdk/_bin/profanity

install-coverage:
	@go install github.com/blend/go-sdk/_bin/coverage

install-recover:
	@go install github.com/blend/go-sdk/_bin/recover
