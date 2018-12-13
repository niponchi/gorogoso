



help:
	@go build
	./goro-goso --help
run:
	@go build -o gorogoso main.go
	./gorogoso -watch=test/**/*.go,runner/**/*.go -entry=test/main.go

build:
	@go build
	mv ./gorogoso ${GOPATH}/bin/gorogoso

unit:
	@go test ./...