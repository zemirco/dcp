
# Discovery and Basic Configuration Protocol (DCP)

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
