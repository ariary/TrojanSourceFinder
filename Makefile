before.build:
	go mod download && go mod vendor

build.tsfinder:
	@go build -o tsfinder cmd/tsfinder/main.go

install.tsfinder:
	@go install github.com/ariary/TrojanSourceFinder/cmd/tsfinder