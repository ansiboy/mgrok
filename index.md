# Ngrok 国内服务

--------------------------

#### Ngrok 搭建

<br/>

#### [点击这里查看如何搭建 ngrok](#build)

-----------------------------------

#### 使用本站提供的 Ngrok 服务

* 下载 ngrok 客户端，请根据使用的操作系统，下载对应的版本

* 运行 ngrok.exe

例如

**建立HTTP隧道**

```
ngrok -subdomain pub -proto=http -config=ngrok.cfg 80
```

**建立数据库隧道**

```
ngrok -proto=tcp -config=ngrok.cfg 3306
```

-----------------------------------

#### 下载

* [windows x64](download/windows_amd64.zip)

* [windows x86](download/windows_386.zip)

* [mac x64](download/darwin_amd64.zip)

* [mac x86](download/darwin_386.zip)

* [linux x64](download/linux_amd64.zip)

* [linux x86](download/linux_386.zip)

* [linux arm](download/linux_arm.zip)