.PHONY:build build-ui build-backend help

build: build-ui build-backend


build-ui:
	cd CasaOS-UI && yarn install && yarn build

build-backend:
	export CGO_ENABLED=1;export CGO_LDFLAGS=-static;go build -o ./casa main.go;upx --lzma --best casa

help:
	@echo "call john"
