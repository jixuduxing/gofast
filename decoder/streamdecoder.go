//streamdecoder

package decoder

import (
	"fmt"
)

const (
	stopbit  = byte(0x80)
	databits = byte(0x7F)
	signbit  = byte(0x40)
)

type streamdecoder struct {
	data []byte
	Pos  int
}

func (sel streamdecoder) datanoprocess() int {
	return len(sel.data) - sel.Pos
}

// value readsuccess
func (sel *streamdecoder) readpmap() ([]byte, bool) {
	beginpos := sel.Pos
	for i := 0; i < 100; i++ {
		if sel.Pos < len(sel.data) {
			if sel.data[sel.Pos]&stopbit > 0 {
				sel.Pos++
				//				fmt.Println("readpmap sel.Pos:", sel.Pos)
				return sel.data[beginpos:sel.Pos], true
			}
			sel.Pos++
		} else {
			fmt.Println("readpmap buffer end")
			// os.Exit(0)
			return []byte{}, false
		}
	}
	fmt.Println("readpmap buffer end  2")
	return []byte{}, false
}

// value index readsuccess
func (sel *streamdecoder) readint() (int, int, bool) {
	var rint int
	var firstch byte
	rint = 0
	for i := 0; i < 5; i++ {
		if sel.Pos < len(sel.data) {
			ch := sel.data[sel.Pos]
			if i == 0 {
				firstch = ch
				if firstch&signbit > 0 {
					rint = -1
				}
			} else if i == 4 {
				if firstch&signbit > 0 {
					if (firstch&databits)>>4 != 7 {
						return 0, 0, false //over flow
					}
				} else if (firstch&databits)>>4 != 0 {
					return 0, 0, false //over flow
				}
			}

			rint <<= 7
			rint |= int(ch & databits)
			sel.Pos++
			if ch&stopbit > 0 {
				return rint, i, true
			}
		}
	}
	return 0, 0, false //no stop bit
}

// value index readsuccess
func (sel *streamdecoder) readint8() (int8, int, bool) {
	var rint int8
	rint = 0
	var firstch byte
	for i := 0; i < 2; i++ {
		if sel.Pos < len(sel.data) {
			ch := sel.data[sel.Pos]
			if i == 0 {
				firstch = sel.data[sel.Pos]
				if firstch&signbit > 0 {
					rint = -1
				}
			} else if i == 1 {
				if firstch&signbit > 0 {
					if (firstch&databits)>>4 != 1 {
						return 0, 0, false //over flow
					}
				} else if (firstch&databits)>>4 != 0 {
					return 0, 0, false //over flow
				}
			}
			rint <<= 7
			rint |= int8(sel.data[sel.Pos] & databits)
			sel.Pos++
			if ch&stopbit > 0 {
				return rint, i, true
			}
		}
	}
	return 0, 0, false //no stop bit
}

// mantissa exponent readsuccess
func (sel *streamdecoder) readdecimal() (int64, int, bool) {
	exponent, i1, flag1 := sel.readint()
	if !flag1 {
		return 0, i1, false
	}
	mantissa, i2, flag2 := sel.readint64()
	if !flag2 {
		return 0, i2, false
	}
	return mantissa, exponent, true
}

// mantissa exponent readsuccess notnull
func (sel *streamdecoder) readdecimalOptional() (int64, int, bool, bool) {
	exponent, i, flag1, flag2 := sel.readintOptional()
	if !flag1 || !flag2 {
		return 0, i, flag1, flag2
	}

	mantissa, _, flag2 := sel.readint64()
	if !flag2 {
		return 0, 0, false, false
	}
	return mantissa, exponent, true, true
}

// value index readsuccess notnull
func (sel *streamdecoder) readuintOptional() (uint, int, bool, bool) {
	if sel.Pos < len(sel.data) {
		ch := sel.data[sel.Pos]
		if ch == stopbit {
			sel.Pos++
			return 0, 0, true, false
		}
	} else {
		return 0, 0, false, false
	}
	rint, _, flag := sel.readuint()
	if !flag {
		return rint, -1, false, false
	}
	if rint > 0 {
		rint--
	}
	return rint, 0, true, true
}

// value index readsuccess notnull

func (sel *streamdecoder) readuint64Optional() (uint64, int, bool, bool) {
	if sel.Pos < len(sel.data) {
		ch := sel.data[sel.Pos]
		if ch == stopbit {
			sel.Pos++
			return 0, 0, true, false // False means NULL
		}
	} else {
		return 0, 0, false, false
	}
	rint, _, flag := sel.readuint64()
	if !flag {
		return rint, -1, false, false
	}
	if rint > 0 {
		rint--
	}
	return rint, 0, true, true
}

// value index readsuccess notnull

func (sel *streamdecoder) readintOptional() (int, int, bool, bool) {
	if sel.Pos < len(sel.data) {
		ch := sel.data[sel.Pos]
		if ch == stopbit {
			sel.Pos++
			return 0, -1, true, false // False means NULL
		}
	} else {
		return 0, -1, false, false
	}
	rint, _, flag := sel.readint()
	if !flag {
		return rint, -1, false, false
	}
	if rint > 0 {
		rint--
	}
	return rint, 0, true, true
}

// value index readsuccess notnull
func (sel *streamdecoder) readint64Optional() (int64, int, bool, bool) {
	if sel.Pos < len(sel.data) {
		ch := sel.data[sel.Pos]
		if ch == stopbit {
			sel.Pos++
			return 0, -1, true, false // means NULL
		}
	} else {
		return 0, -1, false, false
	}
	rint, _, flag := sel.readint64()
	if !flag {
		return rint, 0, false, false
	}
	if rint > 0 {
		rint--
	}
	return rint, 0, true, true
}

// value index readsuccess
func (sel *streamdecoder) readuint() (uint, int, bool) {
	var rint uint
	rint = 0
	var firstch byte
	for i := 0; i < 5; i++ {
		if sel.Pos < len(sel.data) {
			ch := sel.data[sel.Pos]
			if i == 0 {
				firstch = ch

			} else if i == 4 {
				if (firstch&databits)>>4 != 0 {
					return 0, 0, false //over flow
				}
			}
			rint <<= 7
			rint |= uint(sel.data[sel.Pos] & databits)
			sel.Pos++
			if ch&stopbit > 0 {
				return rint, i, true
			}
		} else {
			return 0, 0, false
		}
	}
	return 0, -1, false //no stop bit
}

// value index readsuccess
func (sel *streamdecoder) readuint64() (uint64, int, bool) {
	var rint uint64
	rint = 0
	var firstch byte
	for i := 0; i < 10; i++ {
		if sel.Pos < len(sel.data) {
			ch := sel.data[sel.Pos]
			if i == 0 {
				firstch = ch

			} else if i == 9 {
				if (firstch&databits)>>4 != 0 {
					return 0, 0, false //over flow
				}
			}
			rint <<= 7
			rint |= uint64(sel.data[sel.Pos] & databits)
			sel.Pos++
			if ch&stopbit > 0 {
				return rint, i, true
			}
		} else {
			return 0, 0, false //buffer end
		}
	}
	return 0, -1, false //no stop bit
}

// value index readsuccess
func (sel *streamdecoder) readint64() (int64, int, bool) {
	var rint int64
	rint = 0
	var firstch byte
	for i := 0; i < 10; i++ {
		if sel.Pos < len(sel.data) {
			ch := sel.data[sel.Pos]
			if i == 0 {
				firstch = sel.data[sel.Pos]
				if firstch&signbit > 0 {
					rint = -1
				}
			} else if i == 9 {
				if firstch&signbit > 0 {
					if (firstch&databits)>>4 != 7 {
						return 0, 0, false //over flow
					}
				} else if (firstch&databits)>>4 != 0 {
					return 0, 0, false //over flow
				}
			}
			rint <<= 7
			rint |= int64(sel.data[sel.Pos] & databits)
			sel.Pos++
			if ch&stopbit > 0 {
				return rint, i, true
			}
		} else {
			return 0, 0, false //buffer end
		}
	}
	return 0, -1, false //no stop bit
}

// value index readsuccess
func (sel *streamdecoder) readStringAcii() (string, int, bool) {
	beginpos := sel.Pos
	i := 0
	for true {
		if sel.Pos < len(sel.data) {
			ch := sel.data[sel.Pos]
			sel.Pos++
			if ch&stopbit > 0 {
				return string(append(sel.data[beginpos:sel.Pos-1], ch&databits)), i, true
			}
			i++
		} else {
			return "nil", 0, false
		}
	}
	return "", 0, false
}

// value index readsuccess notnull
func (sel *streamdecoder) readAtringAciiOptional() (string, int, bool, bool) {
	if sel.Pos < len(sel.data) {
		ch := sel.data[sel.Pos]
		if ch == stopbit {
			sel.Pos++
			return "", 0, true, false //FC_NULL_VALUE
		} else if ch == 0 {
			sel.Pos++
			if sel.Pos < len(sel.data) {
				ch = sel.data[sel.Pos]
				if ch == stopbit {
					sel.Pos++
					return "", 0, true, true // #FC_EMPTY_VALUE
				}
				sel.Pos--
			} else {
				return "0", 0, false, false
			}
		}
	} else {
		return "0", 0, false, false
	}
	retstr, i, flag := sel.readStringAcii()
	return retstr, i, flag, true
}

// value index readsuccess notnull
func (sel *streamdecoder) readbyteVectorOptional() ([]byte, int, bool, bool) {
	rlen, i, flag, flag2 := sel.readuint64Optional()
	if !flag || !flag2 || rlen == 0 {
		return []byte{}, i, flag, flag2
	}

	rdata := sel.data[sel.Pos : sel.Pos+int(rlen)]
	sel.Pos += int(rlen)
	return rdata, i + int(rlen), true, true
}

// value index readsuccess
func (sel *streamdecoder) readbyteVector() ([]byte, int, bool) {
	rlen, i, flag := sel.readuint64()
	if !flag || rlen == 0 {
		return []byte{}, i, flag
	}
	rdata := sel.data[sel.Pos : sel.Pos+int(rlen)]
	sel.Pos += int(rlen)
	return rdata, i + int(rlen), true
}
