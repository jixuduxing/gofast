// decoder
package decoder

import (
	"fmt"
)

func generatepmapbits(pmap []byte) {
	fmt.Println("pmap", pmap)
	return
}

func ispresent(seq int, pmap []byte) bool {
	div := seq / 7
	rem := seq % 7

	if div >= len(pmap) {
		return false
	}
	tmp := pmap[div] & (0x40 >> uint(rem))
	if tmp != 0 {
		return true
	}
	return false
}
