// hello
package template

import (
	//	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

const (
	Type_int32      = 0
	Type_uint32     = 1
	Type_int64      = 2
	Type_uint64     = 3
	Type_acsii      = 4
	Type_decimal    = 5
	Type_byteVector = 6
	Type_sequence   = 7
	Type_length     = 8
	Type_typeRef    = 9
	Type_else       = 100
)

const (
	Op_no        = 0
	Op_copy      = 1
	Op_default   = 2
	Op_increment = 3
	Op_constant  = 4
	Op_delta     = 5
	Op_tail      = 6
)

func toitemtype(strtype string) int {
	//	var iret int
	switch {
	case strtype == "string":
		return Type_acsii
	case strtype == "int32":
		return Type_int32
	case strtype == "int64":
		return Type_int64
	case strtype == "uint32":
		return Type_uint32
	case strtype == "length":
		return Type_length
	case strtype == "uint64":
		return Type_uint64
	case strtype == "decimal":
		return Type_decimal
	case strtype == "bytevector":
		return Type_byteVector
	case strtype == "sequence":
		return Type_sequence
	case strtype == "typeRef":
		return Type_typeRef
	default:
		return Type_else
	}
}

func tooptype(strtype string) int {
	switch {
	case strtype == "copy":
		return Op_copy
	case strtype == "default":
		return Op_default
	case strtype == "constant":
		return Op_constant
	case strtype == "increment":
		return Op_increment
	case strtype == "delta":
		return Op_delta
	case strtype == "tail":
		return Op_tail
	default:
		return Op_no
	}
}

type templete struct {
	XMLName xml.Name    `xml:"template"`
	Name    string      `xml:"name,attr"`
	Id      int         `xml:"id,attr"`
	Fields  []tempfield `xml:",any"`
}

type tempfield struct {
	XMLName          xml.Name
	Name             string           `xml:"name,attr"`
	Id               int              `xml:"id,attr"`
	Presence         string           `xml:"presence,attr"`
	DecimalPlaces    string           `xml:"decimalPlaces,attr"`
	Conscontent      tempfieldcontent `xml:"constant"`
	Copycontent      tempfieldcontent `xml:"copy"`
	Defaultcontent   tempfieldcontent `xml:"default"`
	Incrementcontent tempfieldcontent `xml:"increment"`
	Deltacontent     tempfieldcontent `xml:"delta"`
	Tailcontent      tempfieldcontent `xml:"tail"`
	Fields           []tempfield      `xml:",any"`
}

type tempfieldcontent struct {
	XMLName xml.Name
	Value   string `xml:"value,attr"`
}

type tempfile struct {
	XMLName xml.Name `xml:"templates"`
	//	Version    string     `xml:"version,attr"`
	//	UpdateDate string     `xml:"updateDate,attr"`
	//	Xmlns      string     `xml:"xmlns,attr"`
	//	TemplateNs string     `xml:"templateNs,attr"`
	//	Ns         string     `xml:"ns,attr"`
	Templetes []templete `xml:"template"`
}

type Field struct {
	Name        string
	Option      bool   //'false'
	Op          int    //'op_no'
	Datatype    int    //'type_int32'
	Prevalue    string //'0' /*初值*/
	Id          int    //'0'
	Seq         int    //'0'
	Seqlen_item *Field
	Items       []Field
}

type Message struct {
	Msgid   int
	Msgname string
	Fields  []Field
}

type Msgset struct {
	Msgitems map[int]Message
}

func parseElement(token xml.Token) bool {
	//	for t, err = token.Token(); err == nil; t, err = token.Token() {
	//	}
	return true
}

func (self *Field) parsesequence(fld tempfield) bool {
	for _, child := range fld.Fields {
		fiech := Field{Name: child.Name, Id: child.Id, Option: false, Op: Op_no}
		fiech.parseField(child)
		if len(self.Items) == 0 {
			if fiech.Datatype == Type_length {
				self.Seqlen_item = &fiech
				continue
			}
		}
		fiech.Seq = len(self.Items)
		self.Items = append(self.Items, fiech)
	}
	return true
}

func (self *Field) parseField(fld tempfield) bool {
	//	fmt.Println(fld.XMLName.Local, fld.Id, fld.Name, fld.Presence)
	self.Datatype = toitemtype(strings.ToLower(fld.XMLName.Local))
	if self.Datatype == Type_else {
		fmt.Print("err")
		return false
	} else if self.Datatype == Type_typeRef {
		return false
	}
	if fld.Presence == "optional" {
		self.Option = true
	}
	if self.Datatype == Type_sequence {
		return self.parsesequence(fld)
	}

	if len(fld.Conscontent.XMLName.Local) > 0 {
		self.Op = Op_constant
		self.Prevalue = fld.Conscontent.Value
	} else if len(fld.Copycontent.XMLName.Local) > 0 {
		self.Op = Op_copy
		self.Prevalue = fld.Copycontent.Value
	} else if len(fld.Defaultcontent.XMLName.Local) > 0 {
		self.Op = Op_default
		self.Prevalue = fld.Defaultcontent.Value
	} else if len(fld.Incrementcontent.XMLName.Local) > 0 {
		self.Op = Op_increment
		self.Prevalue = fld.Incrementcontent.Value
	} else if len(fld.Deltacontent.XMLName.Local) > 0 {
		self.Op = Op_delta
		self.Prevalue = fld.Deltacontent.Value
	} else if len(fld.Tailcontent.XMLName.Local) > 0 {
		self.Op = Op_tail
		self.Prevalue = fld.Tailcontent.Value
	}
	//	field.op = tooptype(fld.)
	//	fmt.Println("add")
	return true
}

func (self Field) Needplace() bool {
	if self.Op == Op_constant {
		return self.Option
	} else if self.Op == Op_delta {
		return false
	} else if self.Op == Op_no {
		return false
	} else {
		return true
	}
}
func (self *Msgset) ParseTemplate(filename string) bool {
	fmt.Println("ParseTemplate begin!")

	xmlcontent, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return false
	}
	var result tempfile
	err = xml.Unmarshal(xmlcontent, &result)
	if err != nil {
		log.Fatal(err)
		fmt.Println("xml.Unmarshal fail")
		return false
	}
	//	fmt.Println(result)

	self.Msgitems = make(map[int]Message)
	for _, tem := range result.Templetes {
		//		fmt.Println(k, tem.XMLName.Local, tem.Name, tem.Id)
		msg := Message{Msgid: tem.Id, Msgname: tem.Name}
		for _, fld := range tem.Fields {
			//			fmt.Println(j)
			fie := Field{Name: fld.Name, Id: fld.Id, Option: false, Op: Op_no}
			fie.parseField(fld)
			fie.Seq = len(msg.Fields)
			msg.Fields = append(msg.Fields, fie)
		}
		self.Msgitems[tem.Id] = msg
	}
	fmt.Println("ParseTemplate end!")
	return true
}

func main() {
	fmt.Println("Hello World2!")
	//	var v field

	//	v.Id = 100
	mset := Msgset{}
	//	template.ParseTemplate("C:/Users/gao/PycharmProjects/test/shstep/template.xml")
	mset.ParseTemplate("test.xml")
}
