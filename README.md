# Generic FS Mock

This is a Mock Program for the Generic FS API designed with RSA signing.

## Start

Start the mock server first:

```bash
go run -trimpath main.go -pubkey /path/to/pubkey
```

Test api:

```bash
curl -v http://127.0.0.1:8080/download?token=xxx&e=1624250893&t=1624250893&hash=xxx
```

## API Definition

| Method | URL                                                    | Content-Type     | Body                                                                                                                                             | Response                                                             | Desc                                                                                                                      | Priority |
| ------ | ------------------------------------------------------ | ---------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ | -------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------- | -------- |
| GET    | /download?token=xxx&e=1624250893&t=1624250893&hash=xxx | -                | -                                                                                                                                                | file                                                                 | Download file<br>Sign data: `hash`, `e`, `t`                                                                              | 1        |
| POST   | /upload                                                | form-data        | {<br>&ensp;"token":&nbsp;"resource_uuid=xxx&e=1624250893&t=1624250893&token=xxx",<br>&ensp;"file":&nbsp;`binary`<br>}                            | {<br>&ensp;"code":&nbsp;200,<br>&ensp;"message":&nbsp;"success"<br>} | Upload file<br>Sign data: `uuid`, `e`, `t`                                                                                | 1        |
| POST   | /mkzip                                                 | application/json | {<br>&ensp;"token":&nbsp;"xxx",<br>&ensp;"t":&nbsp;1624250893,<br>&ensp;"target_hash":&nbsp;"xxx",<br>&ensp;"source_hashes":&nbsp;"xxx,xxx"<br>} | {<br>&ensp;"code":&nbsp;200,<br>&ensp;"message":&nbsp;"success"<br>} | Package multi files(source_hashes) into one file(target_hash) and store<br>Sign data: `target_hash`, `source_hashes`, `t` | 2        |
| POST   | /persist                                               | application/json | {<br>&ensp;"token":&nbsp;"xxx",<br>&ensp;"t":&nbsp;1624250893,<br>&ensp;"hash":&nbsp;"xxx"<br>}                                                  | {<br>&ensp;"code":&nbsp;200,<br>&ensp;"message":&nbsp;"success"<br>} | Persist file<br>Sign data: `hash`, `t`                                                                                    | 2        |
| POST   | /preupload                                             | application/json | {<br>&ensp;"token":&nbsp;"xxx",<br>&ensp;"t":&nbsp;1624250893,<br>&ensp;"resource_uuid":&nbsp;"xxx"<br>}                                         | {<br>&ensp;"code":&nbsp;200,<br>&ensp;"message":&nbsp;"success"<br>} | Send file information before upload<br>Sign data: `resource_uuid`, `t`                                                    | 1        |
