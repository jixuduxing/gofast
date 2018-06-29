//fastdecoder

package decoder

import (
	"fmt"
	//	"os"
	"gofast/template"
	"strconv"
	// "../template"
)

type fastdecoder struct {
	curmsgid uint
	msgs     template.Msgset
	pmap     []byte
	seq      int
}
type decimal struct {
	mantissa int64
	exponent int
}

//return value,readsuccess
func read(field *template.Field, decod *streamdecoder, isoption bool) (interface{}, bool) {
	if field.Datatype == template.Type_int32 {
		if isoption {
			// fmt.Println
			ret, _, flag, _ := decod.readintOptional()
			return ret, flag
		}
		// fmt.Println
		ret, _, flag := decod.readint()
		return ret, flag

	} else if field.Datatype == template.Type_uint32 {
		if isoption {
			// fmt.Println
			ret, _, flag, _ := decod.readuintOptional()
			return ret, flag
		}
		// fmt.Println
		ret, _, flag := (decod.readuint())
		return ret, flag
	} else if field.Datatype == template.Type_length {
		if isoption {
			// fmt.Println
			ret, _, flag, _ := (decod.readuintOptional())
			return ret, flag
		}
		// fmt.Println
		ret, _, flag := (decod.readuint())
		return ret, flag
	} else if field.Datatype == template.Type_uint64 {
		if isoption {
			// fmt.Println
			ret, _, flag, _ := (decod.readuint64Optional())
			return ret, flag
		}
		// fmt.Println
		ret, _, flag := (decod.readuint64())
		return ret, flag
	} else if field.Datatype == template.Type_int64 {
		if isoption {
			// fmt.Println
			ret, _, flag, _ := (decod.readint64Optional())
			return ret, flag
		}
		// fmt.Println
		ret, _, flag := (decod.readint64())
		return ret, flag
	} else if field.Datatype == template.Type_acsii {
		if isoption {
			// fmt.Println
			ret, _, flag, _ := (decod.readAtringAciiOptional())
			return ret, flag
		}
		// fmt.Println
		ret, _, flag := (decod.readStringAcii())
		return ret, flag
	} else if field.Datatype == template.Type_decimal {
		if isoption {
			mantissa, exponent, flag, _ := decod.readdecimalOptional()
			return decimal{mantissa: mantissa, exponent: exponent}, flag
			// fmt.Println(exponent, mantissa, flag)
		}
		mantissa, exponent, flag := decod.readdecimal()
		return decimal{mantissa: mantissa, exponent: exponent}, flag
		// fmt.Println(exponent, mantissa, i)
	} else if field.Datatype == template.Type_byteVector {
		if isoption {
			ret, _, flag, _ := (decod.readbyteVectorOptional())
			return ret, flag
		}

		ret, _, flag := (decod.readbyteVector())
		return ret, flag

	}
	return nil, false
}

func (sel *fastdecoder) readsequence(fieldseq *template.Field, decod *streamdecoder) bool {
	if len(fieldseq.Items) == 0 {
		fmt.Println(fieldseq.Name+" wrong", 0)
		return false
	}
	sequencelen := 0
	if fieldseq.Seqlen_item.Needplace() {
		if !ispresent(sel.seq, sel.pmap) {
			if fieldseq.Seqlen_item.Op == template.Op_copy {
				sequencelen, _ = strconv.Atoi(fieldseq.Seqlen_item.Prevalue)
			} else {
				fmt.Println("error ")
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
		//		sequencelen,_ =sel.read(fieldseq.Seqlen_item,decod, fieldseq.Seqlen_item.Option)
	}
	sequncedecod := sequencedecoder{}
	fieldseq.Seqlen_item.Prevalue = strconv.Itoa(sequencelen)
	// fmt.Println("enter sequence", fieldseq.Name, sequencelen)
	return sequncedecod.decode(decod, sequencelen, fieldseq)
}

func (sel *fastdecoder) decodermsg(decord *streamdecoder) bool {
	// fmt.Println("msgid = ", sel.curmsgid)
	curmessage, bfind := sel.msgs.Msgitems[sel.curmsgid]
	if !bfind {
		fmt.Println("data can not parse,unknown msgid", sel.curmsgid)
		return false
	}
	sel.seq = 1
	for idx := range curmessage.Fields {
		if idx == 0 {
			continue
		}
		field := &curmessage.Fields[idx]
		if field.Needplace() {
			if field.Datatype == template.Type_else {
				fmt.Println("decoderdata fail for field:1", "test", field.Id, field.Seq)

				return false
			}
			if !ispresent(sel.seq, sel.pmap) {
				// fmt.Println("no seq|id:", sel.seq, field.Id)
				sel.seq++
				continue
			}
			sel.seq++
			if field.Datatype == template.Type_sequence {
				flag := sel.readsequence(field, decord)
				if !flag {
					fmt.Println("readsequence error", field.Name)
					return false
				}
			} else {
				_, flag := read(field, decord, field.Option)
				if !flag {
					fmt.Println("read error", field.Name, field.Id)
					return false
				}
				// fmt.Println("Id1|value:", field.Id, value)

			}
		} else if field.Op == template.Op_no || field.Op == template.Op_delta {
			if field.Datatype == template.Type_else {
				fmt.Println("decoderdata fail for field2:", field.Id, field.Seq)
				return false
			}
			if field.Datatype == template.Type_sequence {
				flag := sel.readsequence(field, decord)
				if !flag {
					fmt.Println("readsequence error", field.Name)
					return false
				}
			} else {
				_, flag := read(field, decord, field.Option)
				if !flag {
					fmt.Println("read error", field.Name, field.Id)
					return false
				}
				// fmt.Println("Id2|value:", field.Id, value)
			}
		}
	}

	return true
}

func (sel *fastdecoder) decodedata(buff []byte) bool {
	decod := &streamdecoder{data: buff, Pos: 0}
	flag := false
	fastmsgcnt := 0
	for decod.datanoprocess() > 0 {
		sel.pmap, flag = decod.readpmap()
		if !flag {
			return false
		}
		generatepmapbits(sel.pmap)

		if ispresent(0, sel.pmap) {
			msgid, _, flag := decod.readuint()
			if !flag {
				return false
			}
			sel.curmsgid = msgid
		} else {
			// fmt.Println("no msgid")
		}
		if !sel.decodermsg(decod) {
			return false
		}
		fastmsgcnt++
	}
	fmt.Println("decodedata end", fastmsgcnt)
	return true
}
