# Bluzelle-go

[![Build Status](https://travis-ci.org/wlwanpan/bluzelle-go.svg?branch=master)](https://travis-ci.org/bluzelle/swarmDB)
[![GoDoc](https://godoc.org/github.com/wlwanpan/bluzelle-go?status.svg)](https://godoc.org/github.com/wlwanpan/bluzelle-go)
[![Gitter chat](https://img.shields.io/gitter/room/nwjs/nw.js.svg?style=flat-square)](https://gitter.im/bluzelle)

## About bluzelle-go

bluzelle-go is a go client built on top of [WebSocket API](https://bluzelle.github.io/api/#websocket-api) that connect to Bluzelle SwarmDB for basic CRUD operations.

## Getting Started

- Installation
```bash
go get github.com/wlwanpan/bluzelle-go
```

- Import
```go
import "github.com/wlwanpan/bluzelle-go"
```

- Initialize
```go
blz := Connect("127.0.0.1", 51010, "80174b53-2dda-49f1-9d6a-6a780d4")
```

- Create
```go
err := blz.Create("key1", []byte("value1"))
```

- Read
```go
value, err := blz.Read("key1")
```

- Update
```go
err := blz.Update("key1", []byte("value2"))
```

- Remove
```go
err := blz.Remove("key1")
```

- Has
```go
has, err := blz.Has("key1")
```

- Keys
```go
keys := blz.Keys()
```

- Size
```go
size := blz.Size()
```

## Reference

Visit the official bluzelle [documentation](https://bluzelle.github.io/api/)
