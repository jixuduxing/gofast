//fastdecoder
package decoder

import (
	"fmt"
	"os"
	"strconv"

	"../template"
)

type fastdecoder struct {
	curmsgid int
	msgs     template.Msgset
	pmap     []byte
	seq      int
}

type sequencedecoder struct {
	seq  int
	pmap []byte
}

func read(field *template.Field, decod *streamdecoder, isoption bool) {
	if field.Datatype == template.Type_int32 {
		if isoption {
			fmt.Println(decod.readint_optional())
			return
		}
		fmt.Println(decod.readint())
		return
	} else if field.Datatype == template.Type_uint32 {
		if isoption {
			fmt.Println(decod.readuint_optional())
			return
		}
		fmt.Println(decod.readuint())
	} else if field.Datatype == template.Type_length {
		if isoption {
			fmt.Println(decod.readuint_optional())
			return
		}
		fmt.Println(decod.readuint())
	} else if field.Datatype == template.Type_uint64 {
		if isoption {
			fmt.Println(decod.readuint64_optional())
			return
		}
		fmt.Println(decod.readuint64())
		return
	} else if field.Datatype == template.Type_int64 {
		if isoption {
			fmt.Println(decod.readint64_optional())
			return
		}
		fmt.Println(decod.readint64())
		return
	} else if field.Datatype == template.Type_acsii {
		if isoption {
			fmt.Println(decod.read_string_acii_optional())
			return
		}
		fmt.Println(decod.read_string_acii())
		return
	} else if field.Datatype == template.Type_decimal {
		if isoption {
			exponent, mantissa, flag := decod.readdecimal_optional()
			fmt.Println(exponent, mantissa, flag)
			return
		}
		exponent, mantissa, i := decod.readdecimal()
		fmt.Println(exponent, mantissa, i)
		return
	} else if field.Datatype == template.Type_byteVector {
		if isoption {
			fmt.Println(decod.readbyteVector_optional())
			return
		}
		fmt.Println(decod.readbyteVector())
		return
	}
	return
}

func (self *fastdecoder) readsequence(fieldseq *template.Field, decod *streamdecoder) {
	if len(fieldseq.Items) == 0 {
		fmt.Println(fieldseq.Name+" wrong", 0)
		return
	}
	sequencelen := 0
	if fieldseq.Seqlen_item.Needplace() {
		if !ispresent(self.seq, self.pmap) {
			if fieldseq.Seqlen_item.Op == template.Op_copy {
				sequencelen, _ = strconv.Atoi(fieldseq.Seqlen_item.Prevalue)
			} else {
				fmt.Println("error ")
			}
		} else {
			tmpvalue := uint(0)
			if fieldseq.Seqlen_item.Option {
				tmpvalue, _ = decod.readuint_optional()
			} else {
				tmpvalue, _ = decod.readuint()
			}
			sequencelen = int(tmpvalue)
			//			self.read(fieldseq.Seqlen_item,decod, fieldseq.Seqlen_item.Option)
		}
		self.seq += 1
	} else {
		tmpvalue := uint(0)
		if fieldseq.Seqlen_item.Option {
			tmpvalue, _ = decod.readuint_optional()
		} else {
			tmpvalue, _ = decod.readuint()
		}
		sequencelen = int(tmpvalue)
		//		sequencelen,_ =self.read(fieldseq.Seqlen_item,decod, fieldseq.Seqlen_item.Option)
	}
	sequncedecod := sequencedecoder{}
	fieldseq.Seqlen_item.Prevalue = strconv.Itoa(sequencelen)
	fmt.Println("enter sequence", fieldseq.Name, sequencelen)
	sequncedecod.decode(decod, sequencelen, fieldseq)
}

func (self *fastdecoder) decodermsg(decord *streamdecoder) bool {
	fmt.Println("msgid = ", self.curmsgid)
	curmessage, bfind := self.msgs.Msgitems[self.curmsgid]
	if !bfind {
		fmt.Println("data can not parse,unknown msgid", self.curmsgid)
		return false
	}
	self.seq = 1
	for idx, _ := range curmessage.Fields[1:] {
		field := &curmessage.Fields[idx+1]
		if field.Needplace() {
			if field.Datatype == template.Type_else {
				fmt.Println("decoderdata fail for field:1", "test", field.Id, field.Seq)

				return false
			}
			if !ispresent(self.seq, self.pmap) {
				fmt.Println("no seq|id:", self.seq, field.Id)
				self.seq += 1
				continue
			}
			self.seq += 1
			fmt.Println("ID|seq:", field.Id, self.seq)
			if field.Datatype == template.Type_sequence {
				self.readsequence(field, decord)
			} else {
				read(field, decord, field.Option)
			}
		} else if field.Op == template.Op_no || field.Op == template.Op_delta {
			if field.Datatype == template.Type_else {
				fmt.Println("decoderdata fail for field2:", field.Id, field.Seq)
				return false
			}
			fmt.Println("data2 ID:", field.Id)
			if field.Datatype == template.Type_sequence {
				self.readsequence(field, decord)
			} else {
				read(field, decord, field.Option)
			}
		}
	}

	return true
}

func (self *fastdecoder) decodedata(buff []byte) bool {
	decod := &streamdecoder{data: buff, Pos: 0}
	for decod.datanoprocess() > 0 {
		self.pmap = decod.readpmap()
		generatepmapbits(self.pmap)

		if ispresent(0, self.pmap) {
			msgid, _ := decod.readint()
			self.curmsgid = msgid
		} else {
			fmt.Println("no msgid")
		}
		if !self.decodermsg(decod) {
			return false
		}
	}
	return true
}

func (self *sequencedecoder) readsequence(fieldseq *template.Field, decod *streamdecoder) {
	if len(fieldseq.Items) == 0 {
		fmt.Println(fieldseq.Name+" wrong", 0)
		return
	}
	sequencelen := 0
	if fieldseq.Seqlen_item.Needplace() {
		if !ispresent(self.seq, self.pmap) {
			if fieldseq.Seqlen_item.Op == template.Op_copy {
				sequencelen, _ = strconv.Atoi(fieldseq.Seqlen_item.Prevalue)
			} else {
				fmt.Println("error ")
				os.Exit(0)
			}
		} else {
			tmpvalue := uint(0)
			if fieldseq.Seqlen_item.Option {
				tmpvalue, _ = decod.readuint_optional()
			} else {
				tmpvalue, _ = decod.readuint()
			}
			sequencelen = int(tmpvalue)
			//			self.read(fieldseq.Seqlen_item,decod, fieldseq.Seqlen_item.Option)
		}
		self.seq += 1
	} else {
		tmpvalue := uint(0)
		if fieldseq.Seqlen_item.Option {
			tmpvalue, _ = decod.readuint_optional()
		} else {
			tmpvalue, _ = decod.readuint()
		}
		sequencelen = int(tmpvalue)
		//		sequencelen,_ =self.read(fieldseq.Seqlen_item,decod, fieldseq.Seqlen_item.Option)
	}
	sequncedecod := sequencedecoder{}
	fieldseq.Seqlen_item.Prevalue = strconv.Itoa(sequencelen)
	fmt.Println("enter sequence", fieldseq.Name, sequencelen)
	sequncedecod.decode(decod, sequencelen, fieldseq)
}

func (self *sequencedecoder) decode(decod *streamdecoder, sequencelen int, fieldseq *template.Field) {
	if sequencelen < 0 {
		os.Exit(0)
	}
	for i := 0; i < sequencelen; i++ {
		self.pmap = decod.readpmap()
		generatepmapbits(self.pmap)
		self.seq = 0
		for idx, _ := range fieldseq.Items {
			field := &fieldseq.Items[idx]
			//			fmt.Println(field.Id, field.Needplace(), field.Seq)
			if field.Needplace() {
				if field.Datatype == template.Type_else {
					fmt.Println("decoderdata fail for field:1", "test", field.Id, field.Seq)
					return
				}
				if !ispresent(self.seq, self.pmap) {
					self.seq += 1
					continue
				}
				self.seq += 1
				fmt.Println("id|seq", field.Id, self.seq)
				if field.Datatype == template.Type_sequence {
					self.readsequence(field, decod)
				} else {
					read(field, decod, field.Option)
				}

			} else if field.Op == template.Op_no || field.Op == template.Op_delta {
				if field.Datatype == template.Type_else {
					fmt.Println("decoderdata fail for field2:", field.Id, field.Seq)
					return
				}
				fmt.Println("data1 ID:", field.Id)
				if field.Datatype == template.Type_sequence {
					self.readsequence(field, decod)
				} else {
					read(field, decod, field.Option)
				}

			}
		}
	}
	fmt.Println("leave sequence ", fieldseq.Name)
}
