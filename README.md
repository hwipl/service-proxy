# service-proxy

TODO

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
