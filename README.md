# Simple HTTP Proxy

This project provides a simple HTTP proxy sever capable of forwarding requests based on the client's IP address. If the client's IP is in the whitelist, the proxy server will forward the requests to the downstream server. Otherwise, it returns a 404 error.

The server can support both individual IP addresses and CIDR notations in the whitelist. It also supports both IPv4 and IPv6.

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
  go build
  ./main
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