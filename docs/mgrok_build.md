

## assets_path 修改

* src/ngrok/client/assets/assets_debug.go 中的 assets_path 

## 代码的编译

在终端，路径为项目的文件夹，例如：

```
d:\projects\mgrok
```

### Ngro 的编译

```
go build -o bin/ngrok.exe -tags "debug"  src/ngrok/main/ngrok/ngrok.go
```

### Ngrod 的编译

```
* go build -o bin/ngrokd.exe -tags "debug"  src/ngrok/main/ngrokd/ngrokd.go
```

### 其他平台版本的编译

#### Linux 64位
```
set GOOS=linux

set GOARCH=amd64

go build -o bin/linux_64/ngrok.exe -tags "debug"  src/ngrok/main/ngrok/ngrok.go

go build -o bin/linux_64/ngrokd.exe -tags "debug"  src/ngrok/main/ngrokd/ngrokd.go
```

#### Linux 32位
```
set GOOS=linux

set GOARCH=386

go build -o bin/linux_32/ngrok.exe -tags "debug"  src/ngrok/main/ngrok/ngrok.go

go build -o bin/linux_32/ngrokd.exe -tags "debug"  src/ngrok/main/ngrokd/ngrokd.go
```

#### Windows 32位
```
set GOOS=windows

set GOARCH=386

go build -o bin/windows_32/ngrok.exe -tags "debug"  src/ngrok/main/ngrok/ngrok.go

go build -o bin/windows_32/ngrokd.exe -tags "debug"  src/ngrok/main/ngrokd/ngrokd.go
```




