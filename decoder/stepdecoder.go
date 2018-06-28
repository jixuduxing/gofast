//stepdecoder BalaBala

package decoder

import (
	"fmt"
	"strconv"

	"gofast/template"
)

//const BalaBala
const (
	RawDataID       = "96"
	BodyLengthID    = "9"
	RawDataLengthID = "95"
	Soh             = '\x01'
)

//Stepdecoder  balabala
type Stepdecoder struct {
	//	curmsgid    int
	Msgs template.Msgset
	//	curtemplate []byte
}

//Decodedata  BalaBala
func (sel Stepdecoder) Decodedata(data []byte) int {
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
		if v == Soh {
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
			if v == Soh {
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

		if id == BodyLengthID {
			bodylen, _ := strconv.Atoi(value)
			fmt.Println(id, "=", value)
			if pos+bodylen < len(data) {
				sel.decodebody(data[pos : pos+bodylen])
				pos += bodylen
			} else {
				return 0
			}
		} else if id == "10" {
			return pos
		}
	}
	return 0
}

func (sel Stepdecoder) decodebody(data []byte) int {
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
		if id != RawDataID {
			SOHfind = false
			for i, v := range data[pos:] {
				if v == Soh {
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
			if id == RawDataLengthID {
				rawdatasize, _ = strconv.Atoi(value)
			}
		} else {
			if rawdatasize > 0 {
				fmt.Println("fast fast fast fast fast fast fast fast ")
				if string(data[pos:pos+2]) == "8=" {
					sel.Decodedata(data[pos : pos+rawdatasize])
					break
				}
				fastdecod := fastdecoder{msgs: sel.Msgs}
				if !fastdecod.decodedata(data[pos : pos+rawdatasize]) {
					fmt.Println("fastdecod.decodedata fail")
				}
				break

			} else {
				SOHfind = false
				for i, v := range data[pos:] {
					if v == Soh {
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
