# go-sse

Simple code written in Golang for handling server-sent events tasks

### server

```bash
go run *.go
```

### client

```bash
# curl
curl -N http://localhost:8888/stream -H 'x-user-id:1'

# httpie
http --stream http://localhost:8888/stream 'x-user-id:1'
```
