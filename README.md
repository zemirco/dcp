
# Discovery and Basic Configuration Protocol (DCP)

[![GoDoc](https://godoc.org/github.com/zemirco/dcp?status.svg)](https://godoc.org/github.com/zemirco/dcp)
[![CircleCI](https://circleci.com/gh/zemirco/dcp.svg?style=svg)](https://circleci.com/gh/zemirco/dcp)

Native Go implementation.

... work in progress ..

## Usage

You need admin rights to use raw sockets.

```sh
go build main.go && sudo ./main
```

If you don't want to use `sudo` give rights via `setcap`.

```sh
sudo setcap cap_net_raw=ep main
```

## UI

There is a demo application with a web based user interface inside `./ui`.

```sh
cd ui
go build main.go && ./main.go
```

Open http://localhost:8085/ in your browser to see a list of all devices in your network.
