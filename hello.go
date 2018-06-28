// hello

package main

import (
	"fmt"
	"gofast/decoder"
	"gofast/template"
	"os"
	"time"
	// "/decoder"
	// "./template"
)

func teststep(filename string) bool {
	fi, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	// content, err := ioutil.ReadFile(filename)
	// if err != nil {
	// 	fmt.Println("Error opening file: ", err)
	// 	return false
	// }
	mset := template.Msgset{}
	mset.ParseTemplate("template.xml")
	//	mset.ParseTemplate("test.xml")
	cnt := 0
	decd := decoder.Stepdecoder{Msgs: mset}
	bufflen := 10 * 1024
	rcontent := make([]byte, bufflen)
	content := []byte{}
	for true {
		nlen, err := fi.Read(rcontent)
		if err != nil {
			fmt.Println("Read error!")
			return false
		}
		fmt.Println("content len|rem:", nlen, len(content))
		if nlen == bufflen {
			content = append(content, rcontent...)
		} else {
			content = append(content, rcontent[:nlen]...)
		}
		// content := append(bytesremain, rcontent[:nlen])
		pos := 0
		for true {
			iret := decd.Decodedata(content[pos:])
			if iret == -1 {
				fmt.Println("wrong buff !")
				return false
			}
			if iret < 1 {
				// fmt.Println("end!")
				content = content[pos:]
				break
			}
			cnt++
			fmt.Println("parsed iret|pos|cnt", iret, pos, cnt)
			pos += iret

		}
	}
	fmt.Println("end!")
	return true
}

func main() {
	fmt.Println("Hello World!")

	//	mset := template.Msgset{}
	//	template.ParseTemplate("C:/Users/gao/PycharmProjects/test/shstep/template.xml")
	//	mset.ParseTemplate("test.xml")
	begintime := time.Now()
	teststep("fastdata")
	endtime := time.Now()
	fmt.Println(begintime, endtime)
}
