package thriftRPC

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"io"
	"log"
	"reflect"
	"sync"
)

type ThriftDecoder struct {
	r 		io.Reader
	mutex	sync.Mutex
}

func NewThriftDecoder(r io.Reader) *ThriftDecoder{
	return &ThriftDecoder{r:r}
}

func (dec *ThriftDecoder) Decode(v interface{}) error {
	//fmt.Println("Begin decoding...")

	//fmt.Println("+++++++++++++++++++++++++")
	//for i:=1; i<=23; i++ {
	//	buf := make([]byte, 1)
	//	if n, err := dec.r.Read(buf[0:1]); n != 1 || err != nil {
	//		break
	//	}
	//	fmt.Println(buf[0])
	//}
	//fmt.Println("-------------------------\n")
	//
	//return nil

	//dec.mutex.Lock()
	//defer dec.mutex.Unlock()

	//version := readInt32(dec.r)
	//var methodName string
	//var typ uint8 = uint8(version & 255)
	//var seqId int32
	//
	//if version < 0 {
	//	version = version & -65536
	//	if version != -2147418112 {
	//		return errors.New("bad version in response message")
	//	} else {
	//		methodName = readString(dec.r)
	//		seqId = readInt32(dec.r)
	//	}
	//}
	//
	//fmt.Println("version = ", version, ", method = ", methodName, ", type = ", typ, ", seqID = ", seqId)



	respType := reflect.TypeOf(v).Elem()
	fieldNum := respType.NumField()
	fieldIndex := 0

	respValue := reflect.ValueOf(v).Elem()

	//read fields
	for ; fieldIndex < fieldNum;  {
		//fmt.Println("....")
		//read field begin
		fieldType := readByte(dec.r)
		fieldId := int16(0)
		if fieldType != 0 {
			fieldId = readInt16(dec.r)
		}
		//0 -- FieldStop
		if fieldType == 0 {
			break
		}

		//0 -- success
		if fieldId == 0 {
			//8 -- int32
			switch fieldType {
			case thrift.I32:
				result := readInt32(dec.r)
				respValue.Field(fieldIndex).Set(reflect.ValueOf(result))
				break
			default:
				log.Println("not implemented for ", fieldId)
			}
			//if fieldType == thrift.I32 {
			//	result := readInt32(dec.r)
			//	respValue.Field(fieldIndex).Set(reflect.ValueOf(result))
			//}
		}
		fieldIndex++
		//fmt.Println(fieldIndex, " -- ", fieldNum)
		//read field end
	}
	//end read fields
	_ = readByte(dec.r)
	//fmt.Println("-------------------------\n")

	return nil
}

func readByte(r io.Reader) byte {
	var buf [2]byte
	if n, err := r.Read(buf[0:1]); n != 1 || err != nil {
		if n != 1 {
			log.Fatal("failed to get enough bytes")
		} else {
			log.Fatal(err.Error())
		}
	}

	return buf[0]
}


func readInt16(r io.Reader) int16 {
	var buf [2]byte
	if n, err := r.Read(buf[0:2]); n != 2 || err != nil {
		log.Fatal(err.Error())
	}
	result := BytesToInt16(buf[0:])
	return result
}


func readInt32(r io.Reader) int32{
	var buf [4]byte
	if n, err := r.Read(buf[0:4]); n != 4 || err != nil {
		log.Fatal(err.Error())
	}
	result := BytesToInt32(buf[0:])
	return result
}

func readString(r io.Reader) string{
	size := readInt32(r)
	if size < 0 {
		log.Fatal("invalid string size")
	}

	buf := make([]byte, size)
	if n, err := r.Read(buf[0:size]); int32(n) != size && err != nil {
		log.Fatal(err.Error())
	}

	return string(buf[0:size])
}