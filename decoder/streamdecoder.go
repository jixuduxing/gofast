//streamdecoder
package decoder

import (
	"fmt"
	"os"
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

func (self streamdecoder) datanoprocess() int {
	return len(self.data) - self.Pos
}

func (self *streamdecoder) readpmap() []byte {
	beginpos := self.Pos
	for i := 0; i < 100; i++ {
		if self.Pos < len(self.data) {
			if self.data[self.Pos]&stopbit > 0 {
				self.Pos += 1
				//				fmt.Println("readpmap self.Pos:", self.Pos)
				return self.data[beginpos:self.Pos]
			}
			self.Pos += 1
		} else {
			fmt.Println("readpmap buffer end")
			os.Exit(0)
		}
	}
	fmt.Println("readpmap buffer end  2")
	return []byte{}
}

func (self *streamdecoder) readint() (int, int) {
	var rint int
	var firstch byte
	rint = 0
	for i := 0; i < 5; i++ {
		if self.Pos < len(self.data) {
			ch := self.data[self.Pos]
			if i == 0 {
				firstch = ch
				if firstch&signbit > 0 {
					rint = -1
				}
			}
			//lue
			rint <<= 7
			rint |= int(ch & databits)
			self.Pos += 1
			if ch&stopbit > 0 {
				return rint, i
			}
		}
	}
	return 0, -1
}

func (self *streamdecoder) readint8() (int8, int) {
	var rint int8
	rint = 0
	var firstch byte
	for i := 0; i < 2; i++ {
		if self.Pos < len(self.data) {
			ch := self.data[self.Pos]
			if i == 0 {
				firstch = self.data[self.Pos]
				if firstch&signbit > 0 {
					rint = -1
				}
			}
			//lue
			rint <<= 7
			rint |= int8(self.data[self.Pos] & databits)
			self.Pos += 1
			if ch&stopbit > 0 {
				return rint, i
			}
		}
	}
	return 0, -1
}

func (self *streamdecoder) readdecimal() (int, int64, int) {
	exponent, i1 := self.readint()
	mantissa, i2 := self.readint64()
	return exponent, mantissa, i1 + i2
}
func (self *streamdecoder) readdecimal_optional() (int, int64, int) {
	exponent, flag1 := self.readint_optional()
	if flag1 == -1 {
		return 0, 0, -1
	}
	mantissa, _ := self.readint64()

	return exponent, mantissa, -1
}
func (self *streamdecoder) readuint_optional() (uint, int) {
	if self.Pos < len(self.data) {
		ch := self.data[self.Pos]
		if ch == stopbit {
			self.Pos += 1
			return 0, -1 // False means NULL
		}
	} else {
		return 0, -1
	}
	rint, _ := self.readuint()
	if rint > 0 {
		rint -= 1
	}
	return rint, 0
}
func (self *streamdecoder) readuint64_optional() (uint64, int) {
	if self.Pos < len(self.data) {
		ch := self.data[self.Pos]
		if ch == stopbit {
			self.Pos += 1
			return 0, -1 // False means NULL
		}
	} else {
		return 0, -1
	}
	rint, _ := self.readuint64()
	if rint > 0 {
		rint -= 1
	}
	return rint, 0
}
func (self *streamdecoder) readint_optional() (int, int) {
	if self.Pos < len(self.data) {
		ch := self.data[self.Pos]
		if ch == stopbit {
			self.Pos += 1
			return 0, -1 // False means NULL
		}
	} else {
		return 0, -1
	}
	rint, _ := self.readint()
	if rint > 0 {
		rint -= 1
	}
	return rint, 0
}
func (self *streamdecoder) readint64_optional() (int64, int) {
	if self.Pos < len(self.data) {
		ch := self.data[self.Pos]
		if ch == stopbit {
			self.Pos += 1
			return 0, -1 // False means NULL
		}
	} else {
		return 0, -1
	}
	rint, _ := self.readint64()
	if rint > 0 {
		rint -= 1
	}
	return rint, 0
}
func (self *streamdecoder) readuint() (uint, int) {
	var rint uint
	rint = 0
	//	var firstch byte
	for i := 0; i < 5; i++ {
		if self.Pos < len(self.data) {
			ch := self.data[self.Pos]
			//			if i == 0 {
			//				firstch = self.data[self.pos]

			//			}
			//lue
			rint <<= 7
			rint |= uint(self.data[self.Pos] & databits)
			self.Pos += 1
			if ch&stopbit > 0 {
				return rint, i
			}
		}
	}
	return 0, -1
}
func (self *streamdecoder) readuint64() (uint64, int) {
	var rint uint64
	rint = 0
	//	var firstch byte
	for i := 0; i < 10; i++ {
		if self.Pos < len(self.data) {
			ch := self.data[self.Pos]
			//			if i == 0 {
			//				firstch = self.data[self.pos]

			//			}
			//lue
			rint <<= 7
			rint |= uint64(self.data[self.Pos] & databits)
			self.Pos += 1
			if ch&stopbit > 0 {
				return rint, i
			}
		}
	}
	return 0, -1
}
func (self *streamdecoder) readint64() (int64, int) {
	var rint int64
	rint = 0
	var firstch byte
	for i := 0; i < 10; i++ {
		if self.Pos < len(self.data) {
			ch := self.data[self.Pos]
			if i == 0 {
				firstch = self.data[self.Pos]
				if firstch&signbit > 0 {
					rint = -1
				}
			}
			//lue
			rint <<= 7
			rint |= int64(self.data[self.Pos] & databits)
			self.Pos += 1
			if ch&stopbit > 0 {
				return rint, i
			}
		}
	}
	return 0, -1
}
func (self *streamdecoder) read_string_acii() (string, int) {
	beginpos := self.Pos
	i := 0
	for true {
		if self.Pos < len(self.data) {
			ch := self.data[self.Pos]
			self.Pos += 1
			if ch&stopbit > 0 {
				return string(append(self.data[beginpos:self.Pos-1], ch&databits)), i
			}
			i += 1
		} else {
			break
		}
	}
	return "", -1
}
func (self *streamdecoder) read_string_acii_optional() (string, int) {
	if self.Pos < len(self.data) {
		ch := self.data[self.Pos]
		if ch == stopbit {
			self.Pos += 1
			return "", -1 //FC_NULL_VALUE
		} else if ch == 0 {
			self.Pos += 1
			if self.Pos < len(self.data) {
				ch = self.data[self.Pos]
				if ch == stopbit {
					self.Pos += 1
					return "", 0 // #FC_EMPTY_VALUE
				} else {
					self.Pos -= 1
				}
			}
		}
	}
	retstr, _ := self.read_string_acii()
	return retstr, 0
}

func (self *streamdecoder) readbyteVector_optional() ([]byte, int) {
	rlen, flag := self.readuint64_optional()
	if flag == -1 {
		return []byte{}, flag
	}

	rdata := self.data[self.Pos : self.Pos+int(rlen)]
	self.Pos += int(rlen)
	return rdata, 0
}

func (self *streamdecoder) readbyteVector() ([]byte, int) {
	rlen, i := self.readuint64()
	if rlen == 0 {
		return []byte{}, i
	}
	rdata := self.data[self.Pos : self.Pos+int(rlen)]
	self.Pos += int(rlen)
	return rdata, 0
}
