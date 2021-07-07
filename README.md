# Generic FS Mock

This is a Mock Program for the Generic FS API designed with RSA signing.

## Start

> Work dir with hard code: uploaded files saved under /tmp/genericfs/files

Start the mock server first:

```bash
go run -trimpath main.go -pubkey /path/to/pubkey
```

Test api:

```bash
curl -v http://127.0.0.1:8080/download?token=xxx&e=1624250893&t=1624250893&hash=xxx
```
