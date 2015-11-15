// Copyright © 2015 Hylke Visser <htdvisser@gmail.com>
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lorawan

import (
	"bytes"
	"reflect"
	"testing"
)

/* PHYPayload Tests */

type PHYPayloadTest struct {
	structure *PHYPayload
	binary    []byte
}

var (
	phyPayloads = []PHYPayloadTest{
		{&PHYPayload{
			MHDR:        mHdrs[1].structure,
			DataPayload: dataPayloads[0].structure,
			MIC:         []byte{0x4d, 0xea, 0xdb, 0x73},
		}, []byte{0xa0, 0x34, 0x12, 0xcd, 0xab, 0x0, 0x2, 0x56, 0x6, 0x54, 0x54, 0x4e, 0x4d, 0xea, 0xdb, 0x73}},
	}
)

func TestParsePHYPayload(t *testing.T) {
	// TODO: Test for invalid inputs

	for _, c := range phyPayloads {
		got, _ := ParsePHYPayload(c.binary)

		if !reflect.DeepEqual(got.MHDR, c.structure.MHDR) {
			t.Errorf("ParseMHDR(%#v).MHDR\n   got: %#v\n  want: %#v", c.binary, got.MHDR, c.structure.MHDR)
		}
		if !reflect.DeepEqual(got.DataPayload.RawFRMPayload, c.structure.DataPayload.RawFRMPayload) {
			t.Errorf("ParseMHDR(%#v).DataPayload.RawFRMPayload\n   got: %#v\n  want: %#v", c.binary, got.DataPayload.RawFRMPayload, c.structure.DataPayload.RawFRMPayload)
		}
		if !bytes.Equal(got.MIC, c.structure.MIC) {
			t.Errorf("ParseMHDR(%#v).MIC\n   got: %#v\n  want: %#v", c.binary, got.MIC, c.structure.MIC)
		}

	}
}

/* MHDR Tests */

type MHDRTest struct {
	structure *MHDR
	binary    byte
}

var (
	mHdrs = []MHDRTest{
		{&MHDR{MType: macMTypeJoinRequest, Major: macMajorLoRaWANR1}, 0x00},
		{&MHDR{MType: macMTypeConfirmedDataDown, Major: macMajorLoRaWANR1}, 0xA0},
		{&MHDR{MType: macMTypeProprietary, Major: 2}, 0xE2},
	}
	parseOnlyMhdrs = []MHDRTest{
		{&MHDR{MType: macMTypeProprietary, Major: 2}, 0xFE},
		{&MHDR{MType: macMTypeJoinRequest, Major: macMajorLoRaWANR1}, 0x1C},
	}
)

func TestMHDRByte(t *testing.T) {
	for _, c := range mHdrs {
		got := c.structure.Byte()
		if got != c.binary {
			t.Errorf("%#v.Bytes()\n   got: %#v\n  want: %#v", c.structure, got, c.binary)
		}
	}
}

func TestParseMHDR(t *testing.T) {
	// TODO: Test for invalid inputs

	for _, c := range append(mHdrs, parseOnlyMhdrs...) {
		got, _ := ParseMHDR(c.binary)
		if got.Major != c.structure.Major {
			t.Errorf("ParseMHDR(%#v).Major\n   got: %#v\n  want: %#v", c.binary, got.Major, c.structure.Major)
		}
		if got.MType != c.structure.MType {
			t.Errorf("ParseMHDR(%#v).MType\n   got: %#v\n  want: %#v", c.binary, got.MType, c.structure.MType)
		}
	}
}
