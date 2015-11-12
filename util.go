// Copyright © 2015 Hylke Visser <htdvisser@gmail.com>
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lorawan

func boolToByte(b bool) byte {
	if b {
		return 0x1
	}
	return 0x0
}
