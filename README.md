# Preface

Tcp traffic duplicator

![](https://skywind3000.github.io/images/p/misc/2024/tcpdup.png)

## Installation

```bash
go install github.com/skywind3000/tcpdup@latest
```

## Usage

Use the following command:

```bash
tcpdup -listen 0.0.0.0:8080 -target 192.168.1.100:8080 \
    -output 127.0.0.1:8081 -input 127.0.0.1:8082
```

1) Redirect TCP traffic from 0.0.0.0:8080 to 192.168.1.100:8080
2) Copy the outgoing traffic (client->server) to 127.0.0.1:8081
3) Copy the incoming traffic (server->client) to 127.0.0.1:8082
4) Both 'output' and 'input' are optional, you can use only one of them.


## Credit

TODO


