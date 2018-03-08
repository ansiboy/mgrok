DIR=
export GOOS=
export GOARCH=

pkgs:
	go get -u github.com/alecthomas/log4go
	go get -u github.com/gizak/termui
	go get -u github.com/inconshreveable/mousetrap
	go get -u github.com/rcrowley/go-metrics
	go get -u gopkg.in/yaml.v1
	go get -u github.com/kardianos/osext
	go get -u github.com/inconshreveable/go-vhost

# build: pkgs
build:
	go build -o .bin/$(DIR)mgrok main/client/mgrok.go
	go build -o .bin/$(DIR)mgrokd main/server/mgrokd.go
	go build -o .bin/$(DIR)httpProxy main/httpProxy/httpProxy.go

linux64: export GOOS=linux
linux64: export GOARCH=amd64
linux64: PLATFORM=linux_amd64
linux64: build

linux32: export GOOS=linux
linux32: export GOARCH=386
linux32: DIR=linux_386/
linux32: build

arm: export GOOS=linux
arm: export GOARCH=arm
arm: DIR=arm/
arm: build

win64: export GOOS=windows
win64: export GOARCH=amd64
win64: DIR=windows_amd64/
win64: build

win32: export GOOS=windows
win32: export GOARCH=386
win32: DIR=windows_386/
win32: build

darwin64: export GOOS=darwin
darwin64: export GOARCH=amd64
darwin64: DIR=darwin_amd64/
darwin64: build

darwin32: export GOOS=darwin
darwin32: export GOARCH=386
darwin32: DIR=darwin_386/
darwin32: build

default: build