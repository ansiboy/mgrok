# Ngrod 使用教程

----------------------------

### Ngrod 命令行参数

**domain**

域名参数，用户通过域名对主机进行访问

例如：

* Mgrok.cn
* Mgrok.cn:80，

端口如果为空，则为监听的端口。有些时候，访问的端口不一定等于监听的端口，比如通过 Nginx 进行转发，则端口要设置为 Nginx 监听的端口。

**httpAddr**

HTTP 协议监听的地址，可以为空。如果为空，监听所有可以 IP 的 80 端口

例如

* :8081
* 127.0.0.1:8081

**httpsAddr**

HTTPS 协议监听的地址，可以为空。如果为空，监听所有可以 IP 的 443 端口

例如

* :8082
* 127.0.0.1:8082

#### 通过 Nginx 转发请求

如果需要通过 Nginx 转发请求，需要加上 X-Host 头，例如：

```
server {
    listen 80;
    server_name t1.mgrok.cn;
    location / {
        proxy_pass http://192.168.1.19:8081;

        proxy_set_header X-Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}

```

### Ngrok 命令行使用

ngrok 的配置文件是完全可选的非常简单 YAML 格式文件，他可以允许你使用 ngrok 一些更高级的功能，例如：

* 同时运行多个隧道
* 连接到自定义的 ngrok 服务器
* 调整 ngrok 一些很神秘的功能

ngrok 的配置文件默认从 ~/.ngrok 加载。你可以通过 -config 参数重写配置文件的地址

#### 同时运行多个隧道

为了运行多个隧道，你需要在配置文件当中使用 tunnels 参数配置每个隧道。隧道的参数以字典的形式配置在配置文件当中。举个例子，让我们来定义三个不同的隧道。第一个隧道是一个有认证的只转发 https 的隧道。第二个隧道转发我们自己机器的 22 端口以便让我可以通过隧道连接到自己的电脑。最后，我们使用自己的域名创造了一个隧道，我们将要在黑客马拉松中展示这个。

```
tunnels:
  client:
    auth: "user:password"
    proto:
      https: 8080
  ssh:
    proto: 
      tcp: 22
  hacks.inconshreveable.com:
    proto:
      http: 9090
```

通过 ngrok start 命令，我们可以同时运行三个隧道，后面要接上我们要启动的隧道名。

```
ngrok start client ssh hacks.inconshreveable.com
```

终端现在看上去应该是这样的：

```
ngrok

Tunnel Status                 online
Version                       1.3/1.3
Forwarding                    https://client.ngrok.com -> 127.0.0.1:8080
Forwarding                    http://hacks.inconshreveable.com -> 127.0.0.1:9090
Forwarding                    tcp://ngrok.com:44764 -> 127.0.0.1:22
```

#### 隧道设置

每一个隧道都可以设置以下五个参数：proto，subdomain，auth，hostname 以及 remote_port。每一个隧道都必须定义 proto ，因为这定义了协议的类型以及转发的目标。当你在运行 http/https 隧道时， auth 参数是可选的，同样， remote_port 也是可选的，他声明了某个端口将要作为远程服务器转发的端口，请注意这只适用于 TCP 隧道。 ngrok 使用每个隧道的名字做到子域名或者域名，但你可以重写他：

```
tunnels:
  client:
    subdomain: "example"
    auth: "user:password"
    proto:
      https: 8080
```

现在当你运行 ngrok start client 的时候，他将会有这样的效果：example.ngrok.com -> 127.0.0.1:8080。相似的，这对自定义域名同样适用，他可以让你通过别名的方式是隧道名更短。

```
tunnels:
  hacks:
    hostname: "hacks.inconshreveable.com"
    proto:
      http: 9090
```

对于 TCP 隧道，你可以会通过 remote_port 参数来指定一个远程服务器的端口作为映射。如果没有声明，服务器将会给你随机分配一个端口。

```
tunnels:
  ssh:
    remote_port: 60123
    proto:
      tcp: 22
```

#### 其他设置选项

通过在配置文件的顶级配置中声明其他可选的选项，ngrok 的配置文件还可以让你做一些更有趣的事情。举个例子，当你在与 ngrok.com 服务器交互的时候可能需要声明 auth_token 。当你需要改变 ngrok 自带的 web 调试工具所绑定的端口是，你可能需要声明 inspect_addr 。

```
auth_token: abc123
inspect_addr: "0.0.0.0:8888"
tunnels:
  ...
```

#### 连接到自定义的 ngrok 服务器

ngrok 支持连接到其他的 ngrokd 服务器上，即便他们并不托管在 ngrok.com 上。首先，显然你必须正确的配置好你的 ngrokd 服务器。如何配置你自己的 ngrokd 服务器请看这里：运行你自己的 ngrokd 服务器。当你运行了你自己的 ngrokd 服务器，你需要设置两个参数来让 ngrok 安全的连接到你的服务器。首先，你需要设置 server_addr 来支出你服务器的地址。然后你需要设置 trust_host_root_certs 来确保你的 TLS 连接安全。

```
server_addr: "example.com:4443"
trust_host_root_certs: true
tunnels:
  ...
```

#### 在 http 代理下运行

最后，你可以设置 ngrok 在 http 代理下运行，这有时候是很有必要的如果你在一个高度限制的企业网络中时。 ngrok 遵守标准的 Unix 环境变量 http_proxy， 但你也可以通过在配置文件中声明 http_proxy 参数来指定。

```
http_proxy: "http://user:password@10.0.0.1:3128"
tunnels:
  ...
```











