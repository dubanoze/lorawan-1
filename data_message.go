// Copyright © 2015 The Things Network
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lorawan

import (
	"bytes"
	"crypto/aes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/jacobsa/crypto/cmac"
)

/* DataPayload Implementations */

// DataPayload contains the data structure for the MAC Payload of a data message
// See Section 4.3 of the LoRaWan Specification
type DataPayload struct {
	FHDR          *FHDR
	RawFHDR       []byte // Use FHDR.Bytes() instead
	FPort         uint8
	RawFRMPayload []byte
}

// Bytes returns the binary representation of the DataPayload
func (dataPayload *DataPayload) Bytes() []byte {
	dataPayloadbuf := new(bytes.Buffer)
	dataPayloadbuf.Write(dataPayload.FHDR.Bytes())
	binary.Write(dataPayloadbuf, binary.LittleEndian, dataPayload.FPort)
	dataPayloadbuf.Write(dataPayload.RawFRMPayload)
	return dataPayloadbuf.Bytes()
}

// ParseDataPayload parses binary data to a DataPayload
func ParseDataPayload(data []byte) (*DataPayload, error) {
	if len(data) < 7 {
		// MACPayload: at least DevAddr(4), FCtrl(1) and FCnt(2)
		return nil, fmt.Errorf("The MACPayload should be at least 7 bytes")
	}

	fHdr, err := ParseFHDR(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse FHDR: %s", err.Error())
	}

	fHdrLen := 7 + len(fHdr.FOpts)

	dataPayload := &DataPayload{
		FHDR:    fHdr,
		RawFHDR: data[:fHdrLen],
	}

	if len(data) > fHdrLen+1 {
		dataPayload.FPort = data[fHdrLen]
		dataPayload.RawFRMPayload = data[fHdrLen+1:]
	}

	return dataPayload, nil
}

/* FHDR Implementations */

// FHDR contains the data structure of a Frame header
// See Section 4.3.1 of the LoRaWan Specification
type FHDR struct {
	DevAddr  uint32
	RawFCtrl byte // Use FCtrl.Bytes() instead
	FCtrl    *FCtrl
	FCnt     uint16
	FOpts    []byte
}

// Bytes returns the binary representation of the FCtrl
func (fHdr *FHDR) Bytes() []byte {
	fHdrbuf := new(bytes.Buffer)
	binary.Write(fHdrbuf, binary.LittleEndian, fHdr.DevAddr)
	fHdrbuf.WriteByte(fHdr.FCtrl.Byte())
	binary.Write(fHdrbuf, binary.LittleEndian, fHdr.FCnt)
	fHdrbuf.Write(fHdr.FOpts)
	return fHdrbuf.Bytes()
}

// ParseFHDR parses binary data to a FHDR
func ParseFHDR(data []byte) (*FHDR, error) {
	fhdr := &FHDR{
		RawFCtrl: data[4],
	}

	fCtrl := ParseFCtrl(fhdr.RawFCtrl)
	fhdr.FCtrl = fCtrl

	binary.Read(bytes.NewReader(data[0:4]), binary.LittleEndian, &fhdr.DevAddr)
	binary.Read(bytes.NewReader(data[5:7]), binary.LittleEndian, &fhdr.FCnt)

	fOptsLen := int(fhdr.FCtrl.FOptsLen)
	index := 7
	if fOptsLen > 0 {
		if len(data) < index+fOptsLen {
			return fhdr, errors.New("The MACPayload does not contain indicated options")
		}
		fhdr.FOpts = data[index : index+fOptsLen]
		index += fOptsLen
	} else {
		fhdr.FOpts = make([]byte, 0)
	}

	return fhdr, nil
}

/* FCtrl Implementations */

// FCtrl contains the data structure of a Frame control byte
// See Section 4.3.1 of the LoRaWan Specification
type FCtrl struct {
	ADR       bool
	ADRACKReq bool
	ACK       bool
	FPending  bool // for downlink messages only
	FOptsLen  uint8
}

// Byte returns the byte representation of the FCtrl
func (fCtrl *FCtrl) Byte() byte {
	return boolToByte(fCtrl.ADR)<<7 |
		boolToByte(fCtrl.ADRACKReq)<<6 |
		boolToByte(fCtrl.ACK)<<5 |
		boolToByte(fCtrl.FPending)<<4 |
		(fCtrl.FOptsLen & 0xF)
}

// ParseFCtrl parses a byte to a FCtrl
func ParseFCtrl(data byte) *FCtrl {
	return &FCtrl{
		ADR:       ((data & 0x80) >> 7) == 1,
		ADRACKReq: ((data & 0x40) >> 6) == 1,
		ACK:       ((data & 0x20) >> 5) == 1,
		FPending:  ((data & 0x10) >> 4) == 1,
		FOptsLen:  (data & 0xF),
	}
}

/* Other Implementations */

// CryptData encrypts or decrypts the Frame Payload for data messages
// See Section 4.3.3 of the LoRaWan Specification
func CryptData(key []byte, data []byte, downlink bool, devAddr uint32, fCnt uint16) ([]byte, error) {
	numBlocks := int(math.Ceil(float64(len(data)) / 16)) // really?

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("Failed to create AES cipher: %s", err.Error())
	}

	s := new(bytes.Buffer)
	for i := 0; i < numBlocks; i++ {
		ai := new(bytes.Buffer)
		ai.Write([]byte{0x01, 0x00, 0x00, 0x00, 0x00})
		ai.WriteByte(boolToByte(downlink))
		binary.Write(ai, binary.LittleEndian, devAddr)
		binary.Write(ai, binary.LittleEndian, uint32(fCnt))
		ai.WriteByte(0x0)
		ai.WriteByte(byte(i + 1))

		si := make([]byte, block.BlockSize())
		block.Encrypt(si, ai.Bytes())

		s.Write(si)
	}
	seq := s.Bytes()

	result := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		result[i] = data[i] ^ seq[i]
	}

	return result, nil
}

// Crypt runs CryptData for the given DataPayload
func (dataPayload *DataPayload) Crypt(key []byte, downlink bool) ([]byte, error) {
	data, err := CryptData(key, dataPayload.RawFRMPayload, downlink, dataPayload.FHDR.DevAddr, dataPayload.FHDR.FCnt)
	if err != nil {
		return nil, fmt.Errorf("Failed to crypt: %s", err.Error())
	}
	return data, nil
}

// CalculateMIC calculates the Message Integrity Code for a data message
// See Section 4.4 of the LoRaWan Specification
func (dataPayload *DataPayload) CalculateMIC(mhdr *MHDR, nwkSKey []byte) ([]byte, error) {
	// msg = MHDR | FHDR | FPORT | FRMPayload
	msgbuf := new(bytes.Buffer)
	msgbuf.WriteByte(mhdr.Byte())
	msgbuf.Write(dataPayload.RawFHDR)
	msgbuf.WriteByte(dataPayload.FPort)
	msgbuf.Write(dataPayload.RawFRMPayload)
	msg := msgbuf.Bytes()

	blocksbuf := new(bytes.Buffer)

	// B0 =  0x49 | 4x 0x00 | Dir (uplink=0x00/downlink=0x01) | DevAddr | FCnt (4 bytes!) | 0x00 | len(msg)
	blocksbuf.Write([]byte{0x49, 0x0, 0x0, 0x0, 0x0})

	switch mhdr.MType {
	case macMTypeUnconfirmedDataUp,
		macMTypeConfirmedDataUp:
		blocksbuf.WriteByte(0x00)
	case macMTypeUnconfirmedDataDown,
		macMTypeConfirmedDataDown:
		blocksbuf.WriteByte(0x01)
	default:
		return nil, fmt.Errorf("Message direction %#v not is neither up, nor down.", mhdr.MType)
	}

	binary.Write(blocksbuf, binary.LittleEndian, dataPayload.FHDR.DevAddr)
	binary.Write(blocksbuf, binary.LittleEndian, uint32(dataPayload.FHDR.FCnt))
	blocksbuf.WriteByte(0x0)
	blocksbuf.WriteByte(byte(len(msg)))

	// Append msg to B0
	blocksbuf.Write(msg)

	blocks := blocksbuf.Bytes()

	hash, err := cmac.New(nwkSKey)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize CMAC: %s", err.Error())
	}

	_, err = hash.Write(blocks)
	if err != nil {
		return nil, fmt.Errorf("Failed to hash data: %s", err.Error())
	}

	return hash.Sum([]byte{})[0:4], nil
}
