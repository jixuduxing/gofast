//stepdecoder
package decoder

import (
	"fmt"
	"strconv"

	"../template"
)

const (
	RAWDATA_ID       = "96"
	BODY_LENGTH_ID   = "9"
	RAWDATALENGTH_ID = "95"
	SOH              = '\x01'
)

type Stepdecoder struct {
	//	curmsgid    int
	Msgs template.Msgset
	//	curtemplate []byte
}

func (self Stepdecoder) Decodedata(data []byte) int {
	if len(data) < 10 {
		fmt.Println("len(data) < 10")
		return 0
	}
	if string(data[:2]) != "8=" {
		fmt.Println("wrong header")
		return -1
	}
	pos := 2
	bSOHfind := false
	var findpos int
	for i, v := range data[2:] {
		if v == SOH {
			bSOHfind = true
			findpos = i
			break
		}
	}
	if !bSOHfind {
		return 0
	}

	value := string(data[pos : pos+findpos])
	//	fmt.Println(value)
	pos += findpos + 1
	var keyfind bool
	for true {
		keyfind = false
		for i, v := range data[pos:] {
			if v == '=' {
				keyfind = true
				findpos = i
				break
			}
		}
		if !keyfind {
			return 0
		}
		id := string(data[pos : pos+findpos])
		pos += findpos + 1
		bSOHfind = false
		for i, v := range data[pos:] {
			if v == SOH {
				bSOHfind = true
				findpos = i
				break
			}
		}
		if !bSOHfind {
			return 0
		}
		value = string(data[pos : pos+findpos])
		pos += findpos + 1

		if id == BODY_LENGTH_ID {
			bodylen, _ := strconv.Atoi(value)
			fmt.Println(id, "=", value)
			self.decodebody(data[pos : pos+bodylen])
			pos += bodylen
		} else if id == "10" {
			return pos
		}
	}
	return 0
}

func (self Stepdecoder) decodebody(data []byte) int {
	//	fmt.Println(string(data))
	if len(data) < 10 {
		return -1
	}
	pos := 0
	rawdatasize := 0
	var keyfind bool
	var SOHfind bool
	var findpos int
	for true {
		keyfind = false
		for i, v := range data[pos:] {
			if v == '=' {
				keyfind = true
				findpos = i
				break
			}
		}
		if !keyfind {
			return -1
		}
		id := string(data[pos : pos+findpos])
		pos += findpos + 1
		if id != RAWDATA_ID {
			SOHfind = false
			for i, v := range data[pos:] {
				if v == SOH {
					SOHfind = true
					findpos = i
					break
				}
			}
			if !SOHfind {
				return -1
			}
			value := string(data[pos : pos+findpos])
			pos += findpos + 1
			if id == RAWDATALENGTH_ID {
				rawdatasize, _ = strconv.Atoi(value)
			}
		} else {
			if rawdatasize > 0 {
				fmt.Println("fast fast fast fast fast fast fast fast ")
				if string(data[pos:pos+2]) == "8=" {
					self.Decodedata(data[pos : pos+rawdatasize])
					break
				}
				fastdecod := fastdecoder{msgs: self.Msgs}
				fastdecod.decodedata(data[pos : pos+rawdatasize])
				break

			} else {
				SOHfind = false
				for i, v := range data[pos:] {
					if v == SOH {
						SOHfind = true
						findpos = i
					}
				}
				if !SOHfind {
					return -1
				}
				//				value := string(data[pos: pos+findpos])
				pos += findpos + 1
			}
		}
	}
	return 0
}
