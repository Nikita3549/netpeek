# netpeek

Fast TCP port scanner written in Go.

```
target google.com  ports 1-1024

OPEN  53/tcp
OPEN  80/tcp
OPEN  443/tcp

open    closed    elapsed
3       1021      0.52s
```

## Install

```bash
go install github.com/Nikita3549/netpeek/cmd/netpeek@latest
```

## Usage

```bash
netpeek scan --host <target> --port <ports>
```

**Examples**

```bash
# Scan all ports
netpeek scan --host localhost --port all

# Scan a range
netpeek scan --host 192.168.1.1 --port 1-1024

# Scan specific ports
netpeek scan --host example.com --port 80,443,8080

# Custom concurrency and timeout
netpeek scan --host example.com --port all -c 500 -t 1000
```

## Flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--host` | `-H` | required | Target host (IP or domain) |
| `--port` | `-p` | required | Ports: `all`, `80`, `1-1024`, `80,443` |
| `--concurrency` | `-c` | 1024 | Number of concurrent goroutines |
| `--timeout` | `-t` | 500 | Timeout per port in milliseconds |

## How it works

netpeek attempts a TCP connection on each port. If the connection succeeds — the port is open. Ports are scanned concurrently via goroutines controlled by a semaphore, making it fast without overwhelming the target.

## License

MIT