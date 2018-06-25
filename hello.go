// hello
package main

import (
	"fmt"

	"io/ioutil"

	"./decoder"
	"./template"
)

func teststep(filename string) bool {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return false
	}
	mset := template.Msgset{}
	mset.ParseTemplate("C:/Users/gao/PycharmProjects/test/shstep/template.xml")
	//	mset.ParseTemplate("test.xml")

	fmt.Println("content len:", len(content))
	decd := decoder.Stepdecoder{Msgs: mset}
	pos := 0
	cnt := 0
	for i := 0; i < 400000; i++ {
		iret := decd.Decodedata(content[pos:])
		if iret < 1 {
			fmt.Println("end!")
			break
		}
		cnt += 1
		fmt.Println("parsed len=", iret, pos, cnt)
		pos += iret

	}
	return true
}

func main() {
	fmt.Println("Hello World!")

	//	mset := template.Msgset{}
	//	template.ParseTemplate("C:/Users/gao/PycharmProjects/test/shstep/template.xml")
	//	mset.ParseTemplate("test.xml")

	teststep("E:/tmp/shrecord_2018-05-2293845")
}
