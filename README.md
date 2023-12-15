# CartApi Service

This is the CartApi service

Generated with

```
micro new --namespace=go.micro --type=api cart-api
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.api.cart-api
- Type: api
- Alias: cart-api

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend etcd.

```
# install etcd
brew install etcd

# run etcd
etcd
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./cart-api-api
```

Build a docker image
```
make docker
```