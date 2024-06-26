# 简易HTTP代理

该项目提供了一个简单的HTTP代理服务器，能够基于客户端IP地址转发请求。如果客户端的IP在白名单中，则代理服务器会将请求转发给下游服务器。否则，它将返回404错误。

服务器既支持白名单中的单个IP地址也支持CIDR表示法。同时支持IPv4和IPv6。

## 处理X-Forwarded-For头部

代理服务器处理`X-Forwarded-For`头部以确定原始的客户端IP地址。这个头部由代理和负载均衡器用来传递它们代理的连接的实际客户端IP。

在`X-Forwarded-For`头部有多个IP地址的情况下，表明请求经过了多个代理，我们的代理服务器将只评估列表中提供的第一个IP地址。头部的格式是由逗号分隔的IP地址列表，最左侧的IP地址是原始的客户端IP。然后，这个IP将与白名单对比，以决定是否应该将请求转发给下游服务器。

例如，如果头部是`X-Forwarded-For: client, proxy1, proxy2`，代理服务器将只考虑`client`来根据白名单配置决定是否允许访问。

## 白名单检查

当收到请求时，代理服务器会将客户端的IP地址与`config.yaml`文件中指定的IP白名单相比较。如果IP在白名单内，无论是作为单个IP地址还是CIDR块的一部分，请求都会被转发给下游服务器。支持IPv4和IPv6地址。

如果客户端的IP地址没有出现在白名单中，或者`X-Forwarded-For`头部包含了不在白名单中的IP，则服务器将返回HTTP 404错误，表示资源未找到或访问不被允许。

## 如何使用

- 安装 [Go](https://golang.org/doc/install)

- 下载或克隆此代码仓库

- 安装依赖项

  ```shell
  go mod download
  ```

- 配置`config.yaml`文件，根据需要替换服务器地址、下游URL，并添加白名单IP或CIDRs：

  ```yaml
  server:
    addr: :8080

  downstream:
    url: https://www.google.com

  whitelist:
    - 192.168.1.1
    - 192.168.1.2
    - 10.22.81.0/24
    - ::1
  ```

- 构建并运行项目

  ```shell
  go build http-proxy
  ./http-proxy
  ```

- 服务器将在指定的地址启动。尝试从列在白名单的客户端访问该地址，以验证请求是否被正确转发。

## 注意

- 服务器端口和下游服务URL可以通过`config.yaml`文件进行配置。

- IP白名单支持单个IP地址和CIDR块。

- 转发请求时会保留X-Forwarded-For和其他头部。

- 如果客户端的IP不在白名单中，则会返回HTTP 404错误。

## 贡献

欢迎提出拉取请求。

请确保相应地更新测试。