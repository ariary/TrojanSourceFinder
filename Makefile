before.build:
	go mod download && go mod vendor

build.tsFinder:
	@echo "build in ${PWD}";go build -o tsFinder cmd/main.go