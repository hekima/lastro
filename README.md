## lastro: A proxy for ruining connections

Lastro is a dead simple proxy that introduces latency between socket 
reads. Lastro is used for simulating slow connections, for example,
between an application server and a database.

## Installing

You need `go` in order to build lastro. Binary packages are currently
not available.

```
go install github.com/zahpee/lastro
```

## Usage

Lastro only have three parameters:

    %  ./lastro -h
    Usage of ./lastro:
      -bind=":1212": address to bind proxy
      -sleep=0: sleep between sockreads in ms
      -target="localhost:80": target address example localhost:80

For running a proxy in port 8081 for a server in port 8080 is as simple as:

    % ./lastro -bind :8081 -target localhost:8080 -sleep 5

This will introduce 5ms sleeps between sockreads.
