// Copyright © 2015 The Things Network
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lorawan

import "fmt"

const (
	// MType bit field values
	macMTypeJoinRequest         = 0
	macMTypeJoinAccept          = 1
	macMTypeUnconfirmedDataUp   = 2
	macMTypeUnconfirmedDataDown = 3
	macMTypeConfirmedDataUp     = 4
	macMTypeConfirmedDataDown   = 5
	// macMTypeRFU = 6
	macMTypeProprietary = 7

	// Major bit field values
	macMajorLoRaWANR1 = 0
	// macMajorRFU = [1:3]
)

/* PHYPayload Implementations */

// PHYPayload contains the data structure of a payload described in
// Section 4.1 of the LoRaWan Specification
type PHYPayload struct {
	MHDR    *MHDR
	RawMHDR byte // Use MHDR.Bytes() instead
	// TODO: Find a more elegant solution for this:
	DataPayload        *DataPayload        // In case it's a data message
	JoinRequestPayload *JoinRequestPayload // In case it's a join request message
	JoinAcceptPayload  *JoinAcceptPayload  // In case it's a join accept message
	RawMACPayload      []byte
	MIC                []byte
}

// MACPayload returns the MAC Payload for this message type
func (phyPayload *PHYPayload) MACPayload() (MACPayload, error) {
	switch phyPayload.MHDR.MType {
	case macMTypeUnconfirmedDataUp,
		macMTypeUnconfirmedDataDown,
		macMTypeConfirmedDataUp,
		macMTypeConfirmedDataDown:
		return phyPayload.DataPayload, nil
	case macMTypeJoinRequest:
		return phyPayload.JoinRequestPayload, nil
	case macMTypeJoinAccept:
		return phyPayload.JoinAcceptPayload, nil
	default:
		return nil, fmt.Errorf("MType %d not supported", phyPayload.MHDR.MType)
	}
}

// ParsePHYPayload parses binary data to a PHYPayload
func ParsePHYPayload(data []byte) (*PHYPayload, error) {
	if len(data) < 5 {
		// PHYPayload: at least MHDR(1) and MIC(4)
		return nil, fmt.Errorf("The PHYPayload should be at least 5 bytes")
	}

	phyPayload := &PHYPayload{
		RawMHDR:       data[0],
		RawMACPayload: data[1 : len(data)-4],
		MIC:           data[len(data)-4 : len(data)],
	}

	mhdr, mHdrErr := ParseMHDR(data[0])
	if mHdrErr != nil {
		return nil, fmt.Errorf("Failed to parse MHDR: %s", mHdrErr.Error())
	}
	phyPayload.MHDR = mhdr

	if mhdr.Major != macMajorLoRaWANR1 {
		return phyPayload, fmt.Errorf("Major version %d not supported", mhdr.Major)
	}

	// TODO: Find a more elegant solution for this:
	var macPldErr error
	switch mhdr.MType {
	case macMTypeUnconfirmedDataUp,
		macMTypeUnconfirmedDataDown,
		macMTypeConfirmedDataUp,
		macMTypeConfirmedDataDown:
		phyPayload.DataPayload, macPldErr = ParseDataPayload(phyPayload.RawMACPayload)
	case macMTypeJoinRequest:
		phyPayload.JoinRequestPayload, macPldErr = ParseJoinRequestPayload(phyPayload.RawMACPayload)
	case macMTypeJoinAccept:
		phyPayload.JoinAcceptPayload, macPldErr = ParseJoinAcceptPayload(phyPayload.RawMACPayload)
	default:
		return phyPayload, fmt.Errorf("MType %d not supported", mhdr.MType)
	}

	return phyPayload, macPldErr
}

/* MHDR Implementations */

// MHDR contains the data structure of a MAC header described in
// Section 4.2 of the LoRaWan Specification
type MHDR struct {
	MType uint8
	Major uint8
}

// Byte returns the byte representation of the MHDR
func (mhdr *MHDR) Byte() byte {
	return (mhdr.MType << 5) | mhdr.Major
}

// ParseMHDR parses binary data to a MHDR
func ParseMHDR(data byte) (*MHDR, error) {
	mhdr := &MHDR{
		MType: (data & 0xe0) >> 5,
		Major: (data & 0x03),
	}

	return mhdr, nil
}

/* MACPayload Implementations */

// MACPayload represents the payload of data/join request/join accept messages
type MACPayload interface {
	Bytes() []byte
	CalculateMIC(mhdr *MHDR, nwkSKey []byte) ([]byte, error)
}
