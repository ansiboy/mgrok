# 搭建 ngrok 服务实现内网穿透

我们经常会有「把本机开发中的 web 项目给朋友看一下」这种临时需求，为此专门在 VPS 上部署一遍就有点太浪费了。之前我通常是在 ADSL 路由器上配个端口映射让本机服务在外网可以访问，但现在大部分运营商不会轻易让你这么干了。一般小运营商也没有公网 IP，自己的路由器出口还是在局域网内，端口映射这种做法就不管用了。

之前我就想过能否借助拥有公网 IP 的主机中转来实现这种需求，后来发现已经有这样的软件了：ngrok。而且 ngrok 官网本身还提供了公共服务，只需要注册一个帐号，运行它的客户端，就可以快速把内网映射出去。不过这么好的服务，没多久就被墙了~

好在 ngrok 是开源的，我在 VPS 上搭了一套服务自己用，一劳永逸地解决了内网穿透这个难题，这里记录一下过程。（注意：ngrok.com 提供的服务是基于 ngrok 2.0，github 上目前只有 1.0 的源码，二者功能和命令有一些区别，用的时候别搞混了）

## 编译 ngrok

我的 VPS 系统是 Ubuntu 14.04.2 LTS，首先装必要的工具：

```
sudo apt-get install build-essential golang mercurial git
```

获取 ngrok 源码：

```
git clone https://github.com/inconshreveable/ngrok.git ngrok
### 请使用下面的地址，修复了无法访问的包地址
git clone https://github.com/tutumcloud/ngrok.git ngrok
cd ngrok
```

生成并替换源码里默认的证书，注意域名修改为你自己的。（之后编译出来的服务端客户端会基于这个证书来加密通讯，保证了安全性）

```
NGROK_DOMAIN="imququ.com"

openssl genrsa -out base.key 2048
openssl req -new -x509 -nodes -key base.key -days 10000 -subj "/CN=$NGROK_DOMAIN" -out base.pem
openssl genrsa -out server.key 2048
openssl req -new -key server.key -subj "/CN=$NGROK_DOMAIN" -out server.csr
openssl x509 -req -in server.csr -CA base.pem -CAkey base.key -CAcreateserial -days 10000 -out server.crt

cp base.pem assets/client/tls/ngrokroot.crt
```

开始编译：

```
sudo make release-server release-client
```

如果一切正常，ngrok/bin 目录下应该有 ngrok、ngrokd 两个可执行文件。

## 服务端

前面生成的 ngrokd 就是服务端程序了，指定证书、域名和端口启动它（证书就是前面生成的，注意修改域名）：

```
sudo ./bin/ngrokd -tlsKey=server.key -tlsCrt=server.crt -domain="imququ.com" -httpAddr=":8081" -httpsAddr=":8082"
```

到这一步，ngrok 服务已经跑起来了，可以通过屏幕上显示的日志查看更多信息。httpAddr、httpsAddr 分别是 ngrok 用来转发 http、https 服务的端口，可以随意指定。ngrokd 还会开一个 4443 端口用来跟客户端通讯（可通过 -tunnelAddr=":xxx" 指定），如果你配置了 iptables 规则，需要放行这三个端口上的 TCP 协议。

现在，通过 https://imququ.com:8081 和 https://imququ.com:8082 就可以访问到 ngrok 提供的转发服务。为了使用方便，建议把域名泛解析到 VPS 上，这样能方便地使用不同子域转发不同的本地服务。我给 imququ.com 做了泛解析，随便访问一个子域，如：http://pub.imququ.com:8081，可以看到这样一行提示：

```
Tunnel pub.imququ.com:8081 not found
```

这说明万事俱备，只差客户端来连了。

## 客户端

如果要把 linux 上的服务映射出去，客户端就是前面生成的 ngrok 文件。但我用的是 Mac，需要指定环境变量再编一次：

```
sudo GOOS=darwin GOARCH=amd64 make release-server release-client
```

这样在 ngrok/bin 目录下会多出来一个 darwin_amd64 目录，这里的 ngrok 文件就可以拷到 Mac 系统用了。

写一个简单的配置文件，随意命名如 ngrok.cfg：

```
server_addr: imququ.com:4443
trust_host_root_certs: false
```

指定子域、要转发的协议和端口，以及配置文件，运行客户端：

```
./ngrok -subdomain pub -proto=http -config=ngrok.cfg 80
```

不出意外可以看到这样的界面，这说明已经成功连上远端服务了：

![](https://st.imququ.com/i/webp/static/uploads/2015/04/ngrok_client.png.webp)

现在再访问 http://pub.imququ.com:8081，访问到的已经是我本机 80 端口上的服务了。

## 管理界面

上面那张 ngrok 客户端运行界面截图中，有一个 Web Interface 地址，这是 ngrok 提供的监控界面。通过这个界面可以看到远端转发过来的 http 详情，包括完整的 request/response 信息，非常方便。

![](https://st.imququ.com/i/webp/static/uploads/2015/04/ngrok_manager.png.webp)

实际上，由于 ngrok 可以转发 TCP，所以还有很多玩法，原理都一样，这里就不多写了。

## 客户端各个版本编译

 编译Liunx 64位客户端&&服务端

```
cd /usr/local/go/src
GOOS=linux GOARCH=amd64 ./make.bash
cd /usr/local/ngrok/
GOOS=linux GOARCH=amd64 make release-server release-client
```

编译Liunx 32位客户端&&服务端

```
cd /usr/local/go/src
GOOS=linux GOARCH=386 ./make.bash
cd /usr/local/ngrok/
GOOS=linux GOARCH=386 make release-server release-client
```

编译 Mac 64位客户端

```
cd /usr/local/go/src
GOOS=darwin GOARCH=amd64 ./make.bash
cd /usr/local/ngrok/
GOOS=darwin GOARCH=amd64 make release-client
```

编译 Mac 32位客户端

```
cd /usr/local/go/src
GOOS=darwin GOARCH=386 ./make.bash
cd /usr/local/ngrok/
GOOS=darwin GOARCH=386 make release-client
```

编译 Windows 64位客户端

```
cd /usr/local/go/src
GOOS=windows GOARCH=amd64 ./make.bash
cd /usr/local/ngrok/
GOOS=windows GOARCH=amd64 make release-client
```

编译 Windows 32位客户端

```
cd /usr/local/go/src
GOOS=windows GOARCH=386 ./make.bash
cd /usr/local/ngrok/
GOOS=windows GOARCH=386 make release-client
```

编译 ARM位客户端

```
cd /usr/local/go/src
GOOS=linux GOARCH=arm ./make.bash
cd /usr/local/ngrok/
GOOS=linux GOARCH=arm make release-client
```

## Ngrok 客户端命令行使用

### HTTP, HTTPS

```
ngrok -subdomain pub -proto=http -config=ngrok.cfg 80
```

* subdomain 是子域名

* proto 是协议，可以是 http, https, 或者 tcp

* 80 是映射的本地端口号

### TCP

```
ngrok -proto=tcp -config=ngrok.cfg 3306
```

* tcp 是协议

* 3306 是本地的端口号


以上文章内容出自：

* https://imququ.com/post/self-hosted-ngrokd.html
* https://y-yun.top/index.php/archives/253/