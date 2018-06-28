//sequencedecoder

package decoder

import (
	"fmt"
	"strconv"

	"gofast/template"
)

type sequencedecoder struct {
	seq  int
	pmap []byte
}

func (sel *sequencedecoder) readsequence(fieldseq *template.Field, decod *streamdecoder) bool {
	if len(fieldseq.Items) == 0 {
		fmt.Println(fieldseq.Name+" wrong", 0)
		return true
	}
	sequencelen := 0
	if fieldseq.Seqlen_item.Needplace() {
		if !ispresent(sel.seq, sel.pmap) {
			if fieldseq.Seqlen_item.Op == template.Op_copy {
				sequencelen, _ = strconv.Atoi(fieldseq.Seqlen_item.Prevalue)
			} else {
				fmt.Println("error ")
				return false
				// os.Exit(0)
			}
		} else {
			tmpvalue := uint(0)
			flag := false
			if fieldseq.Seqlen_item.Option {
				tmpvalue, _, flag, _ = decod.readuintOptional()
			} else {
				tmpvalue, _, flag = decod.readuint()
			}
			if !flag {
				return false
			}
			sequencelen = int(tmpvalue)
			//			sel.read(fieldseq.Seqlen_item,decod, fieldseq.Seqlen_item.Option)
		}
		sel.seq++
	} else {
		tmpvalue := uint(0)
		flag := false
		if fieldseq.Seqlen_item.Option {
			tmpvalue, _, flag, _ = decod.readuintOptional()
		} else {
			tmpvalue, _, flag = decod.readuint()
		}
		if !flag {
			return false
		}
		sequencelen = int(tmpvalue)
	}
	sequncedecod := sequencedecoder{}
	fieldseq.Seqlen_item.Prevalue = strconv.Itoa(sequencelen)
	// fmt.Println("enter sequence", fieldseq.Name, sequencelen)
	return sequncedecod.decode(decod, sequencelen, fieldseq)
}

func (sel *sequencedecoder) decode(decod *streamdecoder, sequencelen int, fieldseq *template.Field) bool {
	if sequencelen < 0 {
		// os.Exit(0)
		return false
	}
	flag := false
	for i := 0; i < sequencelen; i++ {
		sel.pmap, flag = decod.readpmap()
		if !flag {
			return false
		}
		generatepmapbits(sel.pmap)
		sel.seq = 0
		for idx := range fieldseq.Items {
			field := &fieldseq.Items[idx]
			if field.Needplace() {
				if field.Datatype == template.Type_else {
					fmt.Println("decoderdata fail for field:1", "test", field.Id, field.Seq)
					return false
				}
				if !ispresent(sel.seq, sel.pmap) {
					sel.seq++
					continue
				}
				sel.seq++
				if field.Datatype == template.Type_sequence {
					flag := sel.readsequence(field, decod)
					if !flag {
						fmt.Println("readsequence error", field.Name)
						return false
					}
				} else {
					_, flag := read(field, decod, field.Option)
					if !flag {
						fmt.Println("read error", field.Name, field.Id)
						return false
					}
					// fmt.Println("Id3|value:", field.Id, value)
				}

			} else if field.Op == template.Op_no || field.Op == template.Op_delta {
				if field.Datatype == template.Type_else {
					fmt.Println("decoderdata fail for field2:", field.Id, field.Seq)
					return false
				}
				if field.Datatype == template.Type_sequence {
					flag := sel.readsequence(field, decod)
					if !flag {
						fmt.Println("readsequence error", field.Name)
						return false
					}
				} else {
					_, flag := read(field, decod, field.Option)
					if !flag {
						fmt.Println("read error", field.Name, field.Id)
						return false
					}
					// fmt.Println("Id4|value:", field.Id, value)
				}

			}
		}
	}
	// fmt.Println("leave sequence ", fieldseq.Name)
	return true
}
