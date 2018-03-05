## 生成证书

生成证书，替换掉源代码里的证书文件，路径为

```
assets/client/tls/ngrokroot.crt
assets/server/tls/snakeoil.crt
assets/server/tls/snakeoil.key
```

在 Linux 下生成证书

```
NGROK_DOMAIN="*.mgrok.cn"

openssl genrsa -out ngrokroot.key 2048
openssl req -new -x509 -nodes -key ngrokroot.key -days 10000 -subj "/CN=$NGROK_DOMAIN" -out ngrokroot.crt
openssl genrsa -out snakeoil.key 2048
openssl req -new -key snakeoil.key -subj "/CN=$NGROK_DOMAIN" -out snakeoil.csr
openssl x509 -req -in snakeoil.csr -CA ngrokroot.crt -CAkey ngrokroot.key -CAcreateserial -days 10000 -out snakeoil.crt
```

## 把证书等文件嵌入到源码

首先要安装 go-bindata，运行命令

```
go get -u github.com/jteeuwen/go-bindata/...
```

将证书等文件转换成源码

debug 版

```

go-bindata -nomemcopy -pkg=assets -tags=debug -debug=true -o=src/ngrok/client/assets/assets_debug.go assets/client/...

go-bindata -nomemcopy -pkg=assets -tags=debug -debug=true -o=src/ngrok/server/assets/assets_debug.go assets/server/...

```

release 版

```

go-bindata -nomemcopy -pkg=assets -tags=release -debug=false -o=src/ngrok/client/assets/assets_release.go assets/client/...

go-bindata -nomemcopy -pkg=assets -tags=release  -debug=false -o=src/ngrok/server/assets/assets_release.go assets/server/...

```

## 源码的编译

```
go build -o bin/ngrok -tags "debug"  src/ngrok/main/ngrok/ngrok.go

go build -o bin/ngrokd -tags "debug"  src/ngrok/main/ngrokd/ngrokd.go
```

## 其他平台版本的编译

根据目标平台，设置环境变量，例如要生成 linux 64位版本

在 Linux 设置

```
GOOS=linux

GOARCH=amd64
```

在 windows 下设置

```
set GOOS=linux

set GOARCH=amd64
```

### Linux 64位
```
GOOS=linux

GOARCH=amd64

go build -o bin/linux_64/ngrok -tags "debug"  src/ngrok/main/ngrok/ngrok.go

go build -o bin/linux_64/ngrokd -tags "debug"  src/ngrok/main/ngrokd/ngrokd.go
```

### Linux 32位
```
set GOOS=linux

set GOARCH=386

go build -o bin/linux_32/ngrok -tags "debug"  src/ngrok/main/ngrok/ngrok.go

go build -o bin/linux_32/ngrokd -tags "debug"  src/ngrok/main/ngrokd/ngrokd.go
```

### Windows 32位
```
set GOOS=windows

set GOARCH=386

go build -o bin/windows_32/ngrok.exe -tags "debug"  src/ngrok/main/ngrok/ngrok.go

go build -o bin/windows_32/ngrokd.exe -tags "debug"  src/ngrok/main/ngrokd/ngrokd.go
```





