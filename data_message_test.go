// Copyright © 2015 The Things Network
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lorawan

import (
	"bytes"
	"encoding/base64"
	"reflect"
	"testing"
)

/* DataPayload Tests */

type DataPayloadTest struct {
	structure *DataPayload
	binary    []byte
}

var (
	frmPayload   = []byte{0x54, 0x54, 0x4E}
	dataPayloads = []DataPayloadTest{
		{&DataPayload{FHDR: fHdrs[0].structure, FPort: 6, RawFRMPayload: frmPayload}, append(append(fHdrs[0].binary, 0x06), frmPayload...)},
		{&DataPayload{FHDR: fHdrs[1].structure, FPort: 0, RawFRMPayload: []byte{}}, append(fHdrs[1].binary, 0x00)},
	}
)

func TestDataPayloadBytes(t *testing.T) {
	for _, c := range dataPayloads {
		got := c.structure.Bytes()
		if !bytes.Equal(got, c.binary) {
			t.Errorf("%#v.Bytes()\n   got: %#v\n  want: %#v", c.structure, got, c.binary)
		}
	}
}

func TestParseDataPayload(t *testing.T) {
	for _, c := range dataPayloads {
		got, _ := ParseDataPayload(c.binary)

		if got.FPort != c.structure.FPort {
			t.Errorf("ParseDataPayload(%#v).FPort\n   got: %#v\n  want: %#v", c.binary, got.FPort, c.structure.FPort)
		}
		if !bytes.Equal(got.RawFRMPayload, c.structure.RawFRMPayload) {
			t.Errorf("ParseDataPayload(%#v).RawFRMPayload\n   got: %#v\n  want: %#v", c.binary, got.RawFRMPayload, c.structure.RawFRMPayload)
		}
	}

	_, err1 := ParseDataPayload([]byte{0x00, 0x00})
	if err1 == nil {
		t.Errorf("ParseDataPayload should error on invalid data")
	}

	_, err2 := ParseDataPayload([]byte{0xEF, 0xCD, 0x65, 0x87, fCtrls[1].binary, 0x35, 0x40})
	if err2 == nil {
		t.Errorf("ParseDataPayload should error on invalid data")
	}
}

/* FHDR Tests */

type FHDRTest struct {
	structure *FHDR
	binary    []byte
}

var (
	fOpts = []byte{0x13, 0x12, 0x11}
	fHdrs = []FHDRTest{
		{&FHDR{DevAddr: 2882343476, FCtrl: fCtrls[0].structure, FCnt: 22018, FOpts: []byte{}}, []byte{0x34, 0x12, 0xCD, 0xAB, fCtrls[0].binary, 0x02, 0x56}},
		{&FHDR{DevAddr: 2271596015, FCtrl: fCtrls[1].structure, FCnt: 16437, FOpts: fOpts}, append([]byte{0xEF, 0xCD, 0x65, 0x87, fCtrls[1].binary, 0x35, 0x40}, fOpts...)},
	}
	parseOnlyFHdrs = []FHDRTest{
		{&FHDR{DevAddr: 2271596015, FCtrl: fCtrls[1].structure, FCnt: 16437, FOpts: fOpts}, append(append([]byte{0xEF, 0xCD, 0x65, 0x87, fCtrls[1].binary, 0x35, 0x40}, fOpts...), 0x06)},
	}
)

func TestFHDRBytes(t *testing.T) {
	for _, c := range fHdrs {
		got := c.structure.Bytes()
		if !bytes.Equal(got, c.binary) {
			t.Errorf("%#v.Bytes()\n   got: %#v\n  want: %#v", c.structure, got, c.binary)
		}
	}
}

func TestParseFHDR(t *testing.T) {
	for _, c := range append(fHdrs, parseOnlyFHdrs...) {
		got, _ := ParseFHDR(c.binary)
		if got.DevAddr != c.structure.DevAddr {
			t.Errorf("ParseFHDR(%#v).DevAddr\n   got: %#v\n  want: %#v", c.binary, got.DevAddr, c.structure.DevAddr)
		}
		if !reflect.DeepEqual(got.FCtrl, c.structure.FCtrl) {
			t.Errorf("ParseFHDR(%#v).FCtrl\n   got: %#v\n  want: %#v", c.binary, got.FCtrl, c.structure.FCtrl)
		}
		if got.FCnt != c.structure.FCnt {
			t.Errorf("ParseFHDR(%#v).FCnt\n   got: %#v\n  want: %#v", c.binary, got.FCnt, c.structure.FCnt)
		}
		if !bytes.Equal(got.FOpts, c.structure.FOpts) {
			t.Errorf("ParseFHDR(%#v).FOpts\n   got: %#v\n  want: %#v", c.binary, got.FOpts, c.structure.FOpts)
		}
	}

	_, err := ParseFHDR([]byte{0xEF, 0xCD, 0x65, 0x87, fCtrls[1].binary, 0x35, 0x40})
	if err == nil {
		t.Errorf("ParseFHDR should error on invalid data")
	}
}

/* FCtrl Tests */

type FCtrlTest struct {
	structure *FCtrl
	binary    byte
}

var (
	fCtrls = []FCtrlTest{
		{&FCtrl{ADR: false, ADRACKReq: false, ACK: false, FPending: false, FOptsLen: 0}, 0x00},
		{&FCtrl{ADR: true, ADRACKReq: false, ACK: true, FPending: false, FOptsLen: 3}, 0xA3},
		{&FCtrl{ADR: false, ADRACKReq: true, ACK: false, FPending: true, FOptsLen: 15}, 0x5F},
	}
)

func TestFCtrlByte(t *testing.T) {
	for _, c := range fCtrls {
		got := c.structure.Byte()
		if got != c.binary {
			t.Errorf("%#v.Bytes()\n   got: %#v\n  want: %#v", c.structure, got, c.binary)
		}
	}
}

func TestParseFCtrl(t *testing.T) {
	for _, c := range fCtrls {
		got := ParseFCtrl(c.binary)
		if got.ADR != c.structure.ADR {
			t.Errorf("ParseFCtrl(%#v).ADR\n   got: %#v\n  want: %#v", c.binary, got.ADR, c.structure.ADR)
		}
		if got.ADRACKReq != c.structure.ADRACKReq {
			t.Errorf("ParseFCtrl(%#v).ADRACKReq\n   got: %#v\n  want: %#v", c.binary, got.ADRACKReq, c.structure.ADRACKReq)
		}
		if got.ACK != c.structure.ACK {
			t.Errorf("ParseFCtrl(%#v).ACK\n   got: %#v\n  want: %#v", c.binary, got.ACK, c.structure.ACK)
		}
		if got.FPending != c.structure.FPending {
			t.Errorf("ParseFCtrl(%#v).FPending\n   got: %#v\n  want: %#v", c.binary, got.FPending, c.structure.FPending)
		}
		if got.FOptsLen != c.structure.FOptsLen {
			t.Errorf("ParseFCtrl(%#v).FOptsLen\n   got: %#v\n  want: %#v", c.binary, got.FOptsLen, c.structure.FOptsLen)
		}
	}
}

/* Other Tests */

var (
	key = []byte{0x2B, 0x7E, 0x15, 0x16, 0x28, 0xAE, 0xD2, 0xA6, 0xAB, 0xF7, 0x15, 0x88, 0x09, 0xCF, 0x4F, 0x3C}
)

func TestCryptData(t *testing.T) {
	plaintext, _ := base64.StdEncoding.DecodeString("WW91IGxvb2sgZ29vZCwgTG9yYQ==")
	ciphertext, _ := base64.StdEncoding.DecodeString("7Coyr1VW30ZHtGIMUJ7HPHnkCQ==")

	encrypted, _ := CryptData(key, plaintext, true, 2882400018, 43981)
	if !bytes.Equal(encrypted, ciphertext) {
		t.Errorf("CryptData(%#v)\n   got: %#v\n  want: %#v", plaintext, encrypted, ciphertext)
	}

	decrypted, _ := CryptData(key, ciphertext, true, 2882400018, 43981)
	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("CryptData(%#v)\n   got: %#v\n  want: %#v", ciphertext, decrypted, plaintext)
	}

	// TODO: Add more examples

	_, err := CryptData([]byte{0xEF, 0x65, 0x87}, plaintext, true, 2882400018, 43981)
	if err == nil {
		t.Errorf("CryptData should error on invalid data")
	}
}

func TestDataPayloadCrypt(t *testing.T) {
	plaintext, _ := base64.StdEncoding.DecodeString("WW91IGxvb2sgZ29vZCwgTG9yYQ==")
	ciphertext, _ := base64.StdEncoding.DecodeString("OWvMQw/Hk9bgJctqyXYhyVIJ9Q==")

	// TODO: Add more examples

	plain := &DataPayload{
		FHDR:          fHdrs[0].structure,
		FPort:         6,
		RawFRMPayload: plaintext,
	}

	encrypted, _ := plain.Crypt(key, true)

	if !bytes.Equal(encrypted, ciphertext) {
		t.Errorf("DataPayload.Crypt(%#v)\n   got: %#v\n  want: %#v", plain, encrypted, ciphertext)
	}
}

func TestCalculateMIC(t *testing.T) {
	dataPayload := dataPayloads[0].structure
	mHdr := mHdrs[1].structure
	expected := []byte{0x4d, 0xea, 0xdb, 0x73}

	// TODO: Add more examples

	got, _ := dataPayload.CalculateMIC(mHdr, key)

	if !bytes.Equal(got, expected) {
		t.Errorf("DataPayload.CalculateMIC\n   got: %#v\n  want: %#v", got, expected)
	}
}
