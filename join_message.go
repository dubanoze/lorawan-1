// Copyright © 2015 The Things Network
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lorawan

import "fmt"

type JoinRequestPayload struct {
}

func (joinRequestPayload *JoinRequestPayload) Bytes() []byte {
	return []byte{}
}

func (joinRequestPayload *JoinRequestPayload) CalculateMIC(mhdr *MHDR, nwkSKey []byte) ([]byte, error) {
	return nil, nil
}

func ParseJoinRequestPayload(data []byte) (*JoinRequestPayload, error) {
	return nil, fmt.Errorf("Join messages are not implemented yet.")
}

type JoinAcceptPayload struct {
}

func (joinAcceptPayload *JoinAcceptPayload) Bytes() []byte {
	return []byte{}
}

func (joinAcceptPayload *JoinAcceptPayload) CalculateMIC(mhdr *MHDR, nwkSKey []byte) ([]byte, error) {
	return nil, nil
}

func ParseJoinAcceptPayload(data []byte) (*JoinAcceptPayload, error) {
	return nil, fmt.Errorf("Join messages are not implemented yet.")
}
