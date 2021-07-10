# Generic FS Mock

This is a Mock Program for the Generic FS API designed with RSA signing.

## Quick Start

```bash
make
```

## Generate key pairs

```bash
# Generate private key
openssl genrsa -out private-key.pem 1024

# Generate public key
openssl rsa -in private-key.pem -pubout -out public-key.pem
```

## LICENSE

[MIT](./LICENSE)
