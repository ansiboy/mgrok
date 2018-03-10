# MGROK

MGROK 是基于 ngrok 1.7 版本进行修改，精简了一些功能，使得代码更为简洁易于维护。同时也加入了一些功能。

* 精简的功能有

    * 服务端取消了 HTTPS 的支持，如果需要，可用 NGINX 代理实现
    * 服务端与客户的安全连接
    * 服务端的对访问设备的访问授权认证
    * 客户端的 HTTP 请求的监控

* 增加的功能有

    * 支持 NGINX 的代理
    * 服务端支持集群

MGROK 集群架构

```
------------------------------------- 
|        |        |        |        |
| Mgrokd | Mgrokd | Mgrokd | Mgrokd | 
|        |        |        |        |
-------------------------------------
    |                      |
    |TCP                   | HTTP 请求
    |                      |
    |         ----------------------
    |         |     HTTP PROXY     |
    |         ----------------------
    |                       |
    |                       |
------------------------------------
|               NGINX              |
------------------------------------
                    |
                    |
| --------- -------- --------- ----|
|   Mgrok | Mgrok | Mgrok | Mgrok  | 
| --------- -------- ---------- ---|

```

## 已编译文件下载

[点击这里下载](http://www.mgrok.cn)

## 使用

为了方便测试，可以在本机运行。修改 host 文件，加入两个域名

```
127.0.0.1 t.mgrok.cn
127.0.0.1 pub.t.mgrok.cn
```

LINUX 或 MAC 运行 sudo nano /etc/hosts 命令
WINDOWS 修改 windows/system/drives/hosts 文件

### 服务端单例

1. 运行 mgrokd

```yaml
http_addr: :8081
domain: t.mgrok.cn
tunnel_addr: :4444
```

1. 运行 mgrok

```yaml
server_addr: t.mgrok.cn:4444
tunnels:
    ssh:                     
        remote_port: 40022
        proto:
            tcp: 192.168.1.26:22
    maishu:
        proto:
            http: 192.168.1.19:8080
```

### 服务端集群

1. 运行代理 mgrokp，用于代理 HTTP 服务。默认的配置文件为 mgrokd.yaml。主要参数有
    
    * http_addr：HTTP 服务地址，默认为 127.0.0.1:3762
    * data_addr：与 Mgrok 服务端 (Mgrokd) 通讯的地址

    mgrokd.yaml 配置文件

    ```yaml
    http_addr: 127.0.0.1:3762       
    data_addr: 127.0.0.1:6523
    ```

1. 运行服务端 mgrokd，默认的配置文件为 mgrokd.yaml。主要参数有

    * http_addr：8081，HTTP 服务地址
    * domain：t.mgrok.cn，与 HTTP 服务地址对应的域名
    * http_pulbish_port：3762，用于告知客户端真正的访问端口，在这里，外部的 HTTP 服务端口是代理服务的 HTTP 端口 3762

    mgrokd.yaml 配置文件

    ```yaml
    http_addr: :8081
    domain: t.mgrok.cn
    tunnel_addr: :4444
    data_addr: 127.0.0.1:6523
    http_pulbish_port: 3762
    ```

1. 运行 mgrok

    tunnels 配置，请根据实际情况修改

    ```yaml
    server_addr: t.mgrok.cn:4444
    tunnels:
        ssh:                     
            remote_port: 40022
            proto:
                tcp: 192.168.1.26:22
        maishu:
            proto:
                http: 192.168.1.19:8080
    ```

好了，到了这一步，已经接近完成了，先测试一下，如果一切顺利，接着进行最后一步，配置 NGINX

* 首先配置隧道的连接

    关于 NGINX stream 的配置，不了解的同学，自行搜索

    ```conf
    stream {
        upstream mgrok_frontend {
            server 127.0.0.1:4444;
            # server 127.0.0.1:4445;
        }
        server {
            listen 4443;
            proxy_pass mgrok_frontend;
        }
    }
    ```

    * 然后配置 NGINX 的 HTTP 转发

    其中的 **X-Host** 是必须

    ```
    server {
        listen 80 default_server;
        location / {
                proxy_pass http://127.0.0.1:3762;
                proxy_set_header X-Host $http_host;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_redirect ~^https?:\/\/127.0.0.1:8081(.*) http://${host}$1;
                add_header 'Access-Control-Allow-Origin' '*'                   always;
                add_header 'Access-Control-Allow-Methods' '*'                  always;
                add_header 'Access-Control-Allow-Headers' '*'                  always;
        }
    }
    ```

    * 然后修改 mgrokd.yaml，mgrok.yaml 这连个配置文件

        mgrokd.yaml 配置文件

        把 http_pulbish_port 改为 80

        ```yaml
        http_pulbish_port: 80
        ```

        mgrok.yaml 配置文件

        把 server_addr 改为 4443

        ```yaml
        server_addr: t.mgrok.cn:4443
        ```

至此，NINGX 配置已经完成了，你可以配置更多的 MGROKD 实例，然后通过 NGINX 做均衡负载。


## 源码的编译

1. 下载源代码，并拷贝到 GOPATH 文件夹里面。如果不清楚 GOPATH 目录，请运行 go env 命令查看。
1. 如果是在 LINUX 或者 MAC 环境下编译，可以使用 make 命令行编译。各个操作系统版本编译如下

    ```
    make linux64
    make linux32
    make arm
    make win64
    make win32
    make darwin64
    make darwin32
    ```

1. 如果是在 windows 下编译，进入源码文件夹，运行下面的命令行

    ```
    go build -o .bin/mgrok.exe main/client/mgrok.go
    go build -o .bin/mgrokd.exe main/server/mgrokd.go
    go build -o .bin/mgrokp.exe main/proxy/mgrokp.go
    ```

    然后拷贝下面三个配置文件到 .bin 目录
    ```
    main\client\mgrok.yaml
    main\server\mgrokd.yaml
    main\proxy\mgrokp.yaml
    ```

