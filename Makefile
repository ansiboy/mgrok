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

build: pkgs
build:
	go build -o .bin/$(DIR)mgrok main/client/mgrok.go
	go build -o .bin/$(DIR)mgrokd main/server/mgrokd.go
	go build -o .bin/$(DIR)mgrokp main/proxy/mgrokp.go

build_exe:
	go build -o .bin/$(DIR)mgrok.exe main/client/mgrok.go
	go build -o .bin/$(DIR)mgrokd.exe main/server/mgrokd.go
	go build -o .bin/$(DIR)mgrokp.exe main/mgrokp/mgrokp.go

copy:
	cp main/client/mgrok.yaml .bin/$(DIR)mgrok.yaml	
	cp main/server/mgrokd.yaml .bin/$(DIR)mgrokd.yaml
	cp main/proxy/mgrokp.yaml .bin/$(DIR)mgrokp.yaml	

linux64: export GOOS=linux
linux64: export GOARCH=amd64
linux64: DIR=linux_amd64/
linux64: build
linux64: copy

linux32: export GOOS=linux
linux32: export GOARCH=386
linux32: DIR=linux_386/
linux32: build
linux32: copy

arm: export GOOS=linux
arm: export GOARCH=arm
arm: DIR=arm/
arm: build
arm: copy

win64: export GOOS=windows
win64: export GOARCH=amd64
win64: DIR=windows_amd64/
win64: build_exe
win64: copy

win32: export GOOS=windows
win32: export GOARCH=386
win32: DIR=windows_386/
win32: build_exe
win32: copy

darwin64: export GOOS=darwin
darwin64: export GOARCH=amd64
darwin64: DIR=darwin_amd64/
darwin64: build
darwin64: copy

darwin32: export GOOS=darwin
darwin32: export GOARCH=386
darwin32: DIR=darwin_386/
darwin32: build
darwin32: copy

default: build
default: copy