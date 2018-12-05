



help:
	go build
	./goro-goso --help
run:
	go build -o gorogoso main.go
	./gorogoso -watch=test/*.go -entry=test/main.go

build:
	@go build
	mv ./gorogoso ${GOPATH}/bin/gorogoso