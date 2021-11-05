before.build:
	go mod download && go mod vendor

build.tsfinder:
	@echo "build in ${PWD}";go build -o tsfinder cmd/tsfinder/main.go

install.tsfinder:
	@echo "installing tsfinder..";go install github.com/ariary/TrojanSourceFinder/cmd/tsfinder@latest