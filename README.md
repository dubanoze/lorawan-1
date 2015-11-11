# lorawan

Package lorawan provides support for reading and writing LoRaWAN™ messages.

It is intended to grow into a library that provides everything you need to work with LoRaWAN™ in Go.

## Status

**Work in Progress:**

- [ ] Encoders + Decoders
  - [x] `FCtrl`
  - [x] `FHDR`
  - [x] `MACPayload` for data messages
  - [ ] `MACPayload` for join messages
  - [x] `MHDR`
  - [x] `PHYPayload` (needs testing)
- [ ] Crypto
  - [x] Calculating `MIC` (needs testing)
  - [x] Crypto for `FRMPayload` (needs testing)
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

> TODO

## License

MIT licensed, see LICENSE file.
