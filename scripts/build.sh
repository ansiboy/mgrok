#!/bin/sh

go-bindata -nomemcopy -pkg=assets -tags=debug -debug=true -o=src/ngrok/client/assets/assets_debug.go assets/client/...

go-bindata -nomemcopy -pkg=assets -tags=debug -debug=true -o=src/ngrok/server/assets/assets_debug.go assets/server/...

GOPATH=/home/maishu/projects/mgrok

go build -o bin/ngrok -tags "debug"  src/ngrok/main/ngrok/ngrok.go

go build -o bin/ngrokd -tags "debug"  src/ngrok/main/ngrokd/ngrokd.go

cp src/ngrok/main/ngrok/.ngrok bin
