# Simple HTTP Proxy
[中文介绍](./README_CN.md)

This project provides a simple HTTP proxy sever capable of forwarding requests based on the client's IP address. If the client's IP is in the whitelist, the proxy server will forward the requests to the downstream server. Otherwise, it returns a 404 error.

The server can support both individual IP addresses and CIDR notations in the whitelist. It also supports both IPv4 and IPv6.

## Handling X-Forwarded-For Header

The proxy server handles the `X-Forwarded-For` header to determine the original client IP address. This header is used by proxies and load balancers to pass on the actual client IP of a connection that they are proxying.

In cases where there are multiple IP addresses in the `X-Forwarded-For` header, indicating the request has passed through multiple proxies, our proxy server will evaluate only the first IP address provided in the list. The format of the header is a list of IP addresses separated by commas, where the left-most IP address is the original client IP. This IP is then checked against the whitelist to decide if the request should be forwarded to the downstream server.

For example, if the header is `X-Forwarded-For: client, proxy1, proxy2`, the proxy server will only consider `client` to determine if the access is allowed based on the whitelist configuration.

## Whitelist Checking

When a request is received, the proxy server checks the client's IP address against the IP whitelist specified in the `config.yaml` file. If the IP is within the whitelist, either as an individual IP address or as part of a CIDR block, the request is forwarded to the downstream server. Both IPv4 and IPv6 addresses are supported.

If the client's IP address does not appear in the whitelist or if the `X-Forwarded-For` header contains an IP that is not part of the whitelist, the server will return an HTTP 404 error to indicate that the resource is not found or access is not allowed.

## How to Use

- Install [Go](https://golang.org/doc/install)

- Download or clone this repo

- Install the dependencies

  ```shell
  go mod download
  ```

- Configure the `config.yaml` file, replace the server address, downstream URL, and add whitelist IPs or CIDRs as needed:

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

- Build and run the project

  ```shell
  go build http-proxy
  ./http-proxy
  ```

- The server will start at the specified address. Try to visit the address from a client with whitelisted IP to verify that requests are forwarded properly.

## Notes

- The server port and downstream service URL are configurable through the `config.yaml` file.

- The IP whitelist supports both individual IP addresses and CIDR blocks.

- X-Forwarded-For and other headers are preserved when forwarding requests.

- If the client's IP is not in the whitelist, an HTTP 404 error is returned.

## Contributing

Pull requests are welcome.

Please make sure to update the tests as appropriate.