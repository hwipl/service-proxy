# service-proxy

Service-proxy is a reverse proxy server and client. The server is controlled by
the client: clients connect to the server and register services and, thus, open
ports on the server. The server then forwards traffic between these ports and
the respective client ports.

```
           +--------+                         +--------+
           | Server |<---register services:---| Client |
           |        |    port x <-> u         |        |
           |        |    port y <-> v         |        |
           |        |                         |        |
Peer1 <--->[ port x ]<----forward traffic---->[ port u ]
PeerN <--->[ port y ]<----forward traffic---->[ port v ]
           |        |                         |        |
           +--------+                         +--------+
```

## Usage

You can run `service-proxy` with the following command line arguments:

```
  -allowed-ips IPs
        set comma-separated list of IPs the server accepts
        service registrations from, e.g.:
        127.0.0.1,192.168.1.0/24 (default "0.0.0.0/0")
  -allowed-ports ports
        set comma-separated list of ports the server accepts
        in service registrations, e.g.:
        udp:2048-65000,tcp:8000 (default "udp:1024-65535,tcp:1024-65535")
  -c address
        start client and connect to address; requires -r
  -ca-certs files
        read accepted ca-certificates from comma-separated list of files,
        e.g., cert1.pem,cert2.pem,cert3.pem
  -cert file
        read this host's certificate from file, e.g., cert.pem
  -key file
        read the key of this host's certificate from file, e.g., key.pem
  -r services
        register comma-separated list of services on server,
        e.g., tcp:8000:80,udp:53000:53000
  -s address
        start server (default) and listen on address (default ":32323")
```

On a server, it is recommended to use certificates to authenticate clients (see
`-cert`, `-key`, `-ca-certs`), to restrict the address to listen on and the
addresses to accept connections from (see `-s`, `-allowed-ips`), and to
restrict the ports that can be registered (see `-allowed-ports`).

## Examples

Creating a certificate with IP address (SAN) for the server:

```
server$ openssl req -newkey rsa:4096 -new -nodes -x509 -days 3650 \
        -out server-cert.pem -keyout server-key.pem \
        -subj "/C=US/CN=server" -addext "subjectAltName = IP:192.168.1.1"
```

Creating a certificate for a client:

```
client$ openssl req -newkey rsa:4096 -new -nodes -x509 -days 3650 \
        -out client-cert.pem -keyout client-key.pem \
        -subj "/C=US/CN=client"
```

Running the server listening on IP address `192.168.1.1`, in mTLS mode with
certificate `server-cert.pem` and key `server-key.pem`, treating the client
certificate `client-cert.pem` as a CA certificate, accepting connections from
`192.168.1.0/24`, and allowing TCP ports 32000-42000 for service registrations:

```
server$ service-proxy -s 192.168.1.1 \
        -cert server-cert.pem -key server-key.pem \
        -ca-certs client-cert.pem \
        -allowed-ips 192.168.1.0/24 -allowed-ports tcp:32000-42000
```

Running the client connecting to the server's IP address `192.168.1.1`, in mTLS
mode with certificate `client-cert.pem` and key `client-key.pem`, treating the
certificate `server-cert.pem` as a CA certificate, and registering TCP ports
32000 and 32001 on the server with forwarding to port 32000 and 8080 on the
client:

```
client$ service-proxy -c 192.168.1.1 \
        -cert client-cert.pem -key client-key.pem \
        -ca-certs server-cert.pem \
        -r tcp:32000:32000,tcp:32001:8080
```
