# Simple Go HTTP proxy #

[![Build Status][build-img]][build] [![Godoc][godoc-img]][godoc]

[build]: https://github.com/zamibd/http-proxy/actions/workflows/build.yml
[build-img]: https://github.com/zamibd/http-proxy/actions/workflows/build.yml/badge.svg
[godoc]: https://godoc.org/github.com/zamibd/http-proxy
[godoc-img]: https://godoc.org/github.com/zamibd/http-proxy?status.png

A simple HTTP proxy in Go, no caching.
Supports HTTP and CONNECT method for HTTPS & HTTP/2.

It listens on TCP IPv4/IPv6 at the specified port.

Example:

    $ go run http_proxy.go -listen :8080 [-add-forwarded-for]

Flag `-add-forwarded-for` adds client IP to the `X-Forwarded-For` header.

# Setup #

Install go package, create Debian package, install:

    $ go get -u github.com/zamibd/http-proxy
    $ cd $GOPATH/src/github.com/zamibd/http-proxy


# License #

[Apache License, version 2.0](http://www.apache.org/licenses/LICENSE-2.0).

# Bugs, feature requests, questions #

Create a [new issue](https://github.com/zamibd/http-proxy/issues/new).
