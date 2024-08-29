# Toxy

Toxy is a TCP proxy. It terminates your secure connection and distributes
unsecure traffic to your services.

config.toml

```
hostname = "localhost"
port = 44380
cert_file = "./test/server.rsa.crt"
key_file = "./test/server.rsa.key"
load_balancer = "sequential" # "random"
resolve_interval = 10

[[server.web]]
name = "web1"
hostname = "localhost"
port = 8081

[[server.web]]
name = "web2"
hostname = "localhost"
port = 8082

[[server.web]]
name = "web3"
hostname = "localhost"
port = 8083
```

### TODO

- Detect when certificates change and reload the server rather than having to
  do it manually.
