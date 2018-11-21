



help:
	go build
	./goro-goso --help
run:
	go build
	./goro-goso -watch=test/*.go -entry=test/main.go

build:
	@go build
	mv ./goro-goso ${GOPATH}/bin/goro-goso