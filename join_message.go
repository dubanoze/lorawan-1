// Copyright © 2015 Hylke Visser <htdvisser@gmail.com>
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lorawan

import "fmt"

type JoinRequestPayload struct {
}

func (joinRequestPayload *JoinRequestPayload) Bytes() []byte {
	return []byte{}
}

func ParseJoinRequestPayload(data []byte) (*JoinRequestPayload, error) {
	return nil, fmt.Errorf("Join messages are not implemented yet.")
}

type JoinAcceptPayload struct {
}

func (joinAcceptPayload *JoinAcceptPayload) Bytes() []byte {
	return []byte{}
}

func ParseJoinAcceptPayload(data []byte) (*JoinAcceptPayload, error) {
	return nil, fmt.Errorf("Join messages are not implemented yet.")
}
