# Makefile for vynx

run:
	go run main.go
	
build:
	go build -o ./bin/vynx .

build-linux:
	GOOS=linux GOARCH=amd64 go build -o ./bin/vynx_linux_v$(VERSION)_amd64 .
	GOOS=linux GOARCH=arm64 go build -o ./bin/vynx_linux_v$(VERSION)_arm64 .

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/vynx_macos_v$(VERSION)_amd64 .
	GOOS=darwin GOARCH=arm64 go build -o ./bin/vynx_macos_v$(VERSION)_arm64 .

build-windows:
	GOOS=windows GOARCH=amd64 go build -o ./bin/vynx_windows_v$(VERSION)_amd64.exe .
	GOOS=windows GOARCH=arm64 go build -o ./bin/vynx_windows_v$(VERSION)_arm64.exe .

build-all:
	make build-linux
	make build-mac
	make build-windows

tag:
	git tag -a v$(VERSION) -m "Release version $(VERSION)"
	git push origin v$(VERSION)