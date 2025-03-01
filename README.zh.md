# ![Logo](docs/logos/readme.png)

[![License][1]][2] [![Release][3]][4] [![Telegram][5]][6]

[1]: https://img.shields.io/github/license/tobyxdd/hysteria?style=flat-square

[2]: LICENSE.md

[3]: https://img.shields.io/github/v/release/tobyxdd/hysteria?style=flat-square

[4]: https://github.com/tobyxdd/hysteria/releases

[5]: https://img.shields.io/badge/chat-Telegram-blue?style=flat-square

[6]: https://t.me/hysteria_github

Hysteria 是一个功能丰富的，专为恶劣网络环境进行优化的网络工具（双边加速），比如卫星网络、拥挤的公共 Wi-Fi、在中国连接国外服务器等。
基于修改版的 QUIC 协议。目前有以下模式：（仍在增加中）

- SOCKS5 代理 (TCP & UDP)
- HTTP/HTTPS 代理
- TCP/UDP 转发
- TCP/UDP TPROXY 透明代理 (Linux)
- TUN (Windows 下为 TAP)

## 下载安装

### Windows, Linux, macOS CLI

- 从 https://github.com/tobyxdd/hysteria/releases 下载编译好的版本
  - Linux 分为 `hysteria` (带有 tun 支持) 和 `hysteria-notun` (无 tun 支持) 两个版本。无 tun 支持的版本是静态链接，不依赖系统
    glibc 的。**如果你使用了非标准 Linux 发行版，无法正常执行 `hysteria`，可尝试 `hysteria-notun`**
- 使用 Docker 或 Docker Compose: https://github.com/HyNetwork/hysteria/blob/master/Docker.zh.md
- 使用 Arch Linux AUR: https://aur.archlinux.org/packages/hysteria/
- 自己用 `go build ./cmd` 从源码编译

### OpenWrt LuCI app

- [openwrt-passwall](https://github.com/xiaorouji/openwrt-passwall)

### Android

- [SagerNet](https://github.com/SagerNet/SagerNet) 配合 [hysteria-plugin](https://github.com/SagerNet/SagerNet/releases/tag/hysteria-plugin-0.9.2)

### iOS

- [Shadowrocket](https://apps.apple.com/us/app/shadowrocket/id932747118)

## 快速入门

注意：本节提供的配置只是为了快速上手，可能无法满足你的需求。请到 [高级用法](#高级用法) 中查看所有可用选项与含义。

### 服务器

在目录下建立一个 `config.json`

```json
{
  "listen": ":36712",
  "acme": {
    "domains": [
      "your.domain.com"
    ],
    "email": "hacker@gmail.com"
  },
  "obfs": "fuck me till the daylight",
  "up_mbps": 100,
  "down_mbps": 100
}
```

服务端需要一个 TLS 证书。 你可以让 Hysteria 内置的 ACME 尝试自动从 Let's Encrypt 为你的服务器签发一个证书，也可以自己提供。
证书未必一定要是有效、可信的，但在这种情况下客户端需要进行额外的配置。要使用自己的 TLS 证书，参考这个配置：

```json
{
  "listen": ":36712",
  "cert": "/home/ubuntu/my.crt",
  "key": "/home/ubuntu/my.key",
  "obfs": "fuck me till the daylight",
  "up_mbps": 100,
  "down_mbps": 100
}
```

可选的 `obfs` 选项使用提供的密码对协议进行混淆，这样协议就不容易被检测出是 Hysteria/QUIC，可以用来绕过针对性的 DPI 屏蔽或者 QoS。
如果服务端和客户端的密码不匹配就不能建立连接，因此这也可以作为一个简单的密码验证。对于更高级的验证方案请见下文 `auth`。

`up_mbps` 和 `down_mbps` 限制服务器对每个客户端的最大上传和下载速度。这些也是可选的，如果不需要可以移除。

要启动服务端，只需运行

```
./hysteria-linux-amd64 server
```

如果你的配置文件没有命名为 `config.json` 或在别的路径，请用 `-config` 指定路径：

```
./hysteria-linux-amd64 -config blah.json server
```

### 客户端

和服务器端一样，在程序根目录下建立一个`config.json`。

```json
{
  "server": "example.com:36712",
  "obfs": "fuck me till the daylight",
  "up_mbps": 10,
  "down_mbps": 50,
  "socks5": {
    "listen": "127.0.0.1:1080"
  },
  "http": {
    "listen": "127.0.0.1:8080"
  }
}
```

这个配置同时开了 SOCK5 (支持 TCP & UDP) 代理和 HTTP 代理。Hysteria 还有很多其他模式，请务必前往 [高级用法](#高级用法) 了解一下！
要启用/禁用一个模式，在配置文件中添加/移除对应条目即可。

如果你的服务端证书不是由受信任的 CA 签发的，需要用 `"ca": "/path/to/file.ca"` 指定使用的 CA 或者用 `"insecure": true` 忽略所有
证书错误（不推荐）。

`up_mbps` 和 `down_mbps` 在客户端是必填选项，请根据实际网络情况尽量准确地填写，否则将影响 Hysteria 的使用体验。

有些用户可能会尝试用这个功能转发其他加密代理协议，比如 Shadowsocks。这样虽然可行，但从性能的角度不推荐 - Hysteria 本身就用 TLS，
转发的代理协议也是加密的，再加上如今几乎所有网站都是 HTTPS 了，等于做了三重加密。如果需要代理，建议直接使用代理模式。

## 对比

![Bench](docs/bench/bench.png)

## 高级用法

### 服务器

```json5
{
  "listen": ":36712", // 监听地址
  "protocol": "faketcp", // 留空或 "udp", "wechat-video", "faketcp"
  "acme": {
    "domains": [
      "your.domain.com",
      "another.domain.net"
    ], // ACME 证书域名
    "email": "hacker@gmail.com", // 注册邮箱，可选，推荐
    "disable_http": false, // 禁用 HTTP 验证方式
    "disable_tlsalpn": false, // 禁用 TLS-ALPN 验证方式
    "alt_http_port": 8080, // HTTP 验证方式替代端口
    "alt_tlsalpn_port": 4433 // TLS-ALPN 验证方式替代端口
  },
  "cert": "/home/ubuntu/my_cert.crt", // 证书
  "key": "/home/ubuntu/my_key.crt", // 证书密钥
  "up_mbps": 100, // 单客户端最大上传速度
  "down_mbps": 100, // 单客户端最大下载速度
  "disable_udp": false, // 禁用 UDP 支持
  "acl": "my_list.acl", // 见下文 ACL
  "obfs": "AMOGUS", // 混淆密码
  "auth": { // 验证
    "mode": "password", // 验证模式，暂时只支持 "password" 与 "none"
    "config": {
      "password": "yubiyubi"
    }
  },
  "alpn": "ayaya", // QUIC TLS ALPN
  "prometheus_listen": ":8080", // Prometheus 统计接口监听地址 (在 /metrics)
  "recv_window_conn": 15728640, // QUIC stream receive window
  "recv_window_client": 67108864, // QUIC connection receive window
  "max_conn_client": 4096, // 单客户端最大活跃连接数
  "disable_mtu_discovery": false, // 禁用 MTU 探测 (RFC 8899)
  "ipv6_only": false, // 强制把域名解析成 IPv6 地址
  "resolver": "1.1.1.1:53" // DNS 地址
}
```

#### ACME

目前仅支持 HTTP 与 TLS-ALPN 验证方式，不支持 DNS 验证。对于两种方式请分别确保 TCP 80/443 端口能够被访问。

#### 接入外部验证

如果你是商业代理服务提供商，可以这样把 Hysteria 接入到自己的验证后端：

```json5
{
  // ...
  "auth": {
    "mode": "external",
    "config": {
      "http": "https://api.example.com/auth" // 支持 HTTP 和 HTTPS
    }
  }
}
```

对于上述配置，Hysteria 会把验证请求通过 HTTP POST 发送到 `https://api.example.com/auth`

```json5
{
  "addr": "111.222.111.222:52731",
  "payload": "[BASE64]", // 对应客户端配置的 auth 或 auth_str 字段
  "send": 12500000, // 协商后的服务端最大发送速率 (Bps)
  "recv": 12500000 // 协商后的服务端最大接收速率 (Bps)
}
```

后端必须用 HTTP 200 状态码返回验证结果（即使验证不通过）：

```json5
{
  "ok": false,
  "msg": "No idea who you are"
}
```

`ok` 表示验证是否通过，`msg` 是成功/失败消息。

#### Prometheus 流量统计

通过 `prometheus_listen` 选项可以让 Hysteria 暴露一个 Prometheus HTTP 客户端 endpoint 用来统计流量使用情况。
例如如果配置在 8080 端口，则 API 地址是 `http://example.com:8080/metrics`

```text
hysteria_active_conn{auth="55m95auW5oCq"} 32
hysteria_active_conn{auth="aGFja2VyISE="} 7

hysteria_traffic_downlink_bytes_total{auth="55m95auW5oCq"} 122639
hysteria_traffic_downlink_bytes_total{auth="aGFja2VyISE="} 3.225058e+06

hysteria_traffic_uplink_bytes_total{auth="55m95auW5oCq"} 40710
hysteria_traffic_uplink_bytes_total{auth="aGFja2VyISE="} 37452
```

`auth` 是客户端发来的验证密钥，经过 Base64 编码。

### 客户端

```json5
{
  "server": "example.com:36712", // 服务器地址
  "protocol": "faketcp", // 留空或 "udp", "wechat-video", "faketcp"
  "up_mbps": 10, // 最大上传速度
  "down_mbps": 50, // 最大下载速度
  "socks5": {
    "listen": "127.0.0.1:1080", // SOCKS5 监听地址
    "timeout": 300, // TCP 超时秒数
    "disable_udp": false, // 禁用 UDP 支持
    "user": "me", // SOCKS5 验证用户名
    "password": "lmaolmao" // SOCKS5 验证密码
  },
  "http": {
    "listen": "127.0.0.1:8080", // HTTP 监听地址
    "timeout": 300, // TCP 超时秒数
    "user": "me", // HTTP 验证用户名
    "password": "lmaolmao", // HTTP 验证密码
    "cert": "/home/ubuntu/my_cert.crt", // 证书 (变为 HTTPS 代理)
    "key": "/home/ubuntu/my_key.crt" // 证书密钥 (变为 HTTPS 代理)
  },
  "tun": {
    "name": "tun-hy", // TUN 接口名称
    "timeout": 300, // 超时秒数
    "address": "192.0.2.2", // TUN 接口地址（不适用于 Linux）
    "gateway": "192.0.2.1", // TUN 接口网关（不适用于 Linux）
    "mask": "255.255.255.252", // TUN 接口子网掩码（不适用于 Linux）
    "dns": [ "8.8.8.8", "8.8.4.4" ], // TUN 接口 DNS 服务器（仅适用于 Windows）
    "persist": false // 在程序退出之后保留接口（仅适用于 Linux）
  },
  "relay_tcps": [
    {
      "listen": "127.0.0.1:2222", // TCP 转发监听地址
      "remote": "123.123.123.123:22", // TCP 转发目标地址
      "timeout": 300 // TCP 超时秒数
    },
    {
      "listen": "127.0.0.1:13389", // TCP 转发监听地址
      "remote": "124.124.124.124:3389", // TCP 转发目标地址
      "timeout": 300 // TCP 超时秒数
    }
  ],
  "relay_udps": [
    {
      "listen": "127.0.0.1:5333", // UDP 转发监听地址
      "remote": "8.8.8.8:53", // UDP 转发目标地址
      "timeout": 60 // UDP 超时秒数
    },
    {
      "listen": "127.0.0.1:11080", // UDP 转发监听地址
      "remote": "9.9.9.9.9:1080", // UDP 转发目标地址
      "timeout": 60 // UDP 超时秒数
    }
  ],
  "tproxy_tcp": {
    "listen": "127.0.0.1:9000", // TCP 透明代理监听地址
    "timeout": 300 // TCP 超时秒数
  },
  "tproxy_udp": {
    "listen": "127.0.0.1:9000", // UDP 透明代理监听地址
    "timeout": 60 // UDP 超时秒数
  },
  "acl": "my_list.acl", // 见下文 ACL
  "obfs": "AMOGUS", // 混淆密码
  "auth": "[BASE64]", // Base64 验证密钥
  "auth_str": "yubiyubi", // 字符串验证密钥，和上面的选项二选一
  "alpn": "ayaya", // QUIC TLS ALPN
  "server_name": "real.name.com", // 用于验证服务端证书的 hostname
  "insecure": false, // 忽略一切证书错误 
  "ca": "my.ca", // 自定义 CA
  "recv_window_conn": 15728640, // QUIC stream receive window
  "recv_window": 67108864, // QUIC connection receive window
  "disable_mtu_discovery": false, // 禁用 MTU 探测 (RFC 8899)
  "resolver": "1.1.1.1:53" // DNS 地址
}
```

#### 伪装 TCP (faketcp 模式)

某些网络可能对 UDP 流量施加各种限制，或者完全屏蔽。Hysteria 提供了一个 "faketcp" 模式，让服务端与客户端之间用看起来是 TCP 但实际不走
系统 TCP 栈的方式通信。通过这种方式可以让防火墙、QoS 设备认为这是真的 TCP 连接，绕过对 UDP 的限制。

目前只在 Linux 上支持（客户端和服务器都是），并且需要 root 权限。

如果你的服务器有防火墙，请放行相应的 TCP 端口而不是 UDP。

#### 透明代理

TPROXY 模式 (`tproxy_tcp` 和 `tproxy_udp`) 只在 Linux 下可用。

参考阅读：
- https://www.kernel.org/doc/Documentation/networking/tproxy.txt
- https://powerdns.org/tproxydoc/tproxy.md.html

## 优化建议

### 针对超高传速度进行优化

如果要用 Hysteria 进行极高速度的传输 (如内网超过 10G 或高延迟跨国超过 1G)，请增加系统的 UDP receive buffer 大小。

```shell
sysctl -w net.core.rmem_max=4000000
```

这个命令会在 Linux 下将 buffer 大小提升到 4 MB 左右。

你可能还需要提高 `recv_window_conn` 和 `recv_window` (服务器端是 `recv_window_client`) 以确保它们至少不低于带宽-延迟的乘积。
比如如果想在一条 RTT 200ms 的线路上达到 500 MB/s 的速度，receive window 至少需要 100 MB (500*0.2)

### 路由器与其他嵌入式设备

对于运算性能和内存十分有限的嵌入式设备，如果不是必须的话建议关闭混淆，可以带来少许性能提升。

Hysteria 服务端与客户端默认的 receive window 大小是 64 MB。如果设备内存不够，请考虑通过配置降低。建议保持 stream receive window
和 connection receive window 之间 1:4 的比例关系。

## 关于 ACL

[ACL 文件格式](ACL.zh.md)

ACL 在服务端和客户端都可以使用。在服务端可以用来实现限制客户端能访问的目标，对客户端任何模式都有效。在客户端只有 SOCKS5 和 HTTP 代理
支持 ACL。其他模式下没有效果（所有流量都会走代理）。

## URI Scheme

希望包含链接分享/导入功能的第三方客户端，建议按照如下 URI Scheme 实现（最初由 Shadowrocket 引入）：

    hysteria://host:port?protocol=udp&auth=123456&peer=sni.domain&insecure=1&upmbps=100&downmbps=100&alpn=hysteria&obfs=xplus&obfsParam=123456#remarks

    - host: hostname or IP address of the server to connect to (required)
    - port: port of the server to connect to (required)
    - protocol: protocol to use ("udp", "wechat-video", "faketcp") (optional, default: "udp")
    - auth: authentication payload (string) (optional)
    - peer: SNI for TLS (optional)
    - insecure: ignore certificate errors (optional)
    - upmbps: upstream bandwidth in Mbps (required)
    - downmbps: downstream bandwidth in Mbps (required)
    - alpn: QUIC ALPN (optional)
    - obfs: Obfuscation mode (optional, empty or "xplus")
    - obfsParam: Obfuscation password (optional)
    - remarks: remarks (optional)

## 日志

程序默认在 stdout 输出 DEBUG 级别，文字格式的日志。

如果需要修改日志级别可以使用 `LOGGING_LEVEL` 环境变量，支持 `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace`

如果需要输出 JSON 可以把 `LOGGING_FORMATTER` 设置为 `json`

如果需要修改日志时间戳格式可以使用 `LOGGING_TIMESTAMP_FORMAT`


 ## 自定义 CA 方法

  1. 假设服务器地址是 `123.123.123.123`, 端口`5678`UDP/TCP协议未被防火墙拦截
  2. 已经安装了 openssl
  3. hysteria 已经安装在 `/root/hysteria/`目录下
<details>
  <summary>4. 生成自定义CA证书</summary>

- 在 `/root/hysteria/` 目录下，将以下shell命令保存为 `generate.sh` , 并赋予执行权限: `chmod +x ./generate.sh` 后，运行 `./generate.sh` 命令生成自定义CA证书
- 或者在`/root/hysteria/` 目录下，直接执行以下shell命令生成自定义CA证书

``` shell
#!/usr/bin/env bash

domain=$(openssl rand -hex 8)
password=$(openssl rand -hex 16)
obfs=$(openssl rand -hex 6)
path="/root/hysteria"
# 生成CAkey
openssl genrsa -out hysteria.ca.key 2048
# 生成CA证书
openssl req -new -x509 -days 3650 -key hysteria.ca.key -subj "/C=CN/ST=GD/L=SZ/O=Hysteria, Inc./CN=Hysteria Root CA" -out hysteria.ca.crt

openssl req -newkey rsa:2048 -nodes -keyout hysteria.server.key -subj "/C=CN/ST=GD/L=SZ/O=Hysteria, Inc./CN=*.${domain}.com" -out hysteria.server.csr
# 签发服务端用的证书
openssl x509 -req -extfile <(printf "subjectAltName=DNS:${domain}.com,DNS:www.${domain}.com") -days 3650 -in hysteria.server.csr -CA hysteria.ca.crt -CAkey hysteria.ca.key -CAcreateserial -out hysteria.server.crt

cat > ./client.json <<EOF
{
    "server": "123.123.123.123:5678",
    "alpn": "h3",
    "obfs": "${obfs}",
    "auth_str": "${password}",
    "up_mbps": 30,
    "down_mbps": 30,
    "socks5": {
        "listen": "0.0.0.0:1080"
    },
    "http": {
        "listen": "0.0.0.0:8080"
    },
    "server_name": "www.${domain}.com",
    "ca": "${path}/hysteria.ca.crt"
}
EOF


cat > ./server.json <<EOF
{
    "listen": ":5678",
    "alpn": "h3",
    "obfs": "${obfs}",
    "cert": "${path}/hysteria.server.crt",
    "key": "${path}/hysteria.server.key" ,
    "auth": {
        "mode": "password",
        "config": {
            "password": "${password}"
        }
    }
}
EOF
```
</details>

5. 服务端：复制 `server.json`、 `hysteria.server.crt`、 `hysteria.server.key` 到 `/root/hysteria/` 目录下，运行 `/root/hysteria/hysteria -c /root/hysteria/server.json server` 命令

6. 客户端：假设客户端运行目录也为`/root/hysteria`, 复制 `client.json`、`hysteria.ca.crt` 到 `/root/hysteria/` 目录下，运行 `/root/hysteria/hysteria -c /root/hysteria/client.json` 命令

7. 生成CA证书之后，根据自身情况修改服务器地址、端口和证书文件路径，加上`obfs`和`alpn`是防止首次在某些环境下被墙，第一次在全参数情况下测试通过后，可以自身网络环境删除不必须要参数，比如`obfs`和`alpn`.

8. iOS 端如果使用的是小火箭 Shadowrocket，可以把文件`hysteria.ca.crt` airdrop到手机，然后在手机上安装并信任后, 就可以使用自定义CA证书了。
