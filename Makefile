# Makefile for vynx

OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m)
OS_ARCH := $(OS)_$(ARCH)

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

build-all:
	make build-linux
	make build-mac

tag:
	git tag -a v$(VERSION) -m "Release version $(VERSION)"
	git push origin v$(VERSION)

install:
	@make build
	@echo "Installing vynx to /usr/local/bin/vynx"
	@sudo cp ./bin/vynx /usr/local/bin/vynx

uninstall:
	@rm -f /usr/local/bin/vynx

clean:
	@rm -rf ./bin

gen-seckey:
	KEY=$$(LC_ALL=C tr -dc 'A-Z2-7' </dev/urandom | head -c 32); if grep -q "^AES_KEY=" .env 2>/dev/null; then sed -i "" "s/^AES_KEY=.*/AES_KEY=$$KEY/" .env; else printf "AES_KEY=%s\n" "$$KEY" >> .env; fi;