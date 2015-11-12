# lorawan

Package lorawan provides support for reading and writing LoRaWAN™ messages.

It is intended to grow into a library that provides everything you need to work with LoRaWAN™ in Go.

## Status

[![Build Status](https://travis-ci.org/htdvisser/lorawan.svg)](https://travis-ci.org/htdvisser/lorawan)

**Work in Progress:**

- [ ] Encoders + Decoders
  - [x] `FCtrl`
  - [x] `FHDR`
  - [x] `MACPayload` for data messages
  - [ ] `MACPayload` for join request messages
  - [ ] `MACPayload` for join accept messages
  - [x] `MHDR`
  - [x] `PHYPayload` (works, but needs testing)
- [ ] Crypto
  - [x] Calculating `MIC` (works, but needs testing)
  - [x] Crypto for `FRMPayload` (works, but needs testing)
  - [ ] Crypto for join accept messages
- [ ] Convenience Functions

**For the future:**

- [ ] MAC Commands
- [ ] End Device Activation
- [ ] Class B devices
- [ ] Class C devices

## Install

```
go get github.com/htdvisser/lorawan
```

## Usage

```go
import (
    "github.com/htdvisser/lorawan"
)
```

> TODO

## Docs

[![GoDoc](https://godoc.org/github.com/htdvisser/lorawan?status.svg)](https://godoc.org/github.com/htdvisser/lorawan)

## Links

- [The LoRaWAN™ network protocol (LoRa® Alliance)](http://www.lora-alliance.org/)
- [The Things Network](http://thethingsnetwork.org/)

## License

MIT licensed, see LICENSE file.
