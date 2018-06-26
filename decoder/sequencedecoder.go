//sequencedecoder
package decoder

import (
	"fmt"
	"os"
	"strconv"

	"../template"
)

type sequencedecoder struct {
	seq  int
	pmap []byte
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
				//				fmt.Println("id|seq", field.Id, self.seq)
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
				//				fmt.Println("data1 ID:", field.Id)
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
