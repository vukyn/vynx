# Makefile for vynx

run:
	go run main.go
	
build:
	go build -o ./bin/vynx .

build-linux:
	GOOS=linux GOARCH=amd64 go build -o ./bin/vynx_linux_amd64 .

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/vynx_macos_amd64 .

build-windows:
	GOOS=windows GOARCH=amd64 go build -o ./bin/vynx_amd64.exe .