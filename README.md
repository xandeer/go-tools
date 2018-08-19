# go-tools

## hs

A static http server.

Install:

```shell
go get -u github.com/xandeer/go-tools/hs
```

Usage:

```shell
# hs -d dir -p port
hs -d ~/tmp # default ./
hs -p 1234 # default 9999
```

## ghh

Listen for github webhook, just for push event. Run `make` after pull.

Install:

```shell
go get -u github.com/xandeer/go-tools/ghh
```

Usage:

```shell
# hs -d dir -p port -s secret -b branches
hs -s serect -d ~/tmp # default ./
hs -s serect -p 1234 # default 3001
hs -s serect -b master,dev # default master
```
