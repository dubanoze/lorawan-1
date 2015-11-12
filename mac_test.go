// Copyright © 2015 Hylke Visser <htdvisser@gmail.com>
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lorawan

import "testing"

/* PHYPayload Tests */

func TestPHYPayloadRaw(t *testing.T) {
	// TODO: Implement TestPHYPayloadRaw
}

func TestParsePHYPayload(t *testing.T) {
	// TODO: Implement TestParsePHYPayload
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
