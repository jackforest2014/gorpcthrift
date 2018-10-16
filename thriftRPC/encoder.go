package thriftRPC

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

type ThriftEncoder struct {
	w 				io.Writer
	buffer 			*bytes.Buffer

}

func NewThriftEncoder(r io.Writer) *ThriftEncoder{
	return &ThriftEncoder{w:r}
}

func (enc *ThriftEncoder) Encode(v interface{}) error {
	req := v.(*thriftRequest)
	var msgType int32 = 1
	var version = -2147418112 | msgType

	buf := make([]byte, 1024)
	offset := 0
	offset += Int32ToBytes(version, buf[offset:])
	offset += Int32ToBytes(int32(len(req.Method)), buf[offset:])
	offset += StringToBytes(req.Method, buf[offset:])
	//fmt.Println("req.Id = ", req.Id)
	offset += Int32ToBytes(int32(req.Id), buf[offset:])
	//offset += StringToBytes(req.Method, buf[offset:])

	//fmt.Println("============================")
	//fmt.Println(req.Param)
	param := reflect.ValueOf(req.Param).Elem()
	//fmt.Println(param.Type())

	paramType := reflect.TypeOf(reflect.ValueOf(req.Param).Elem())
	//paramValue := reflect.ValueOf(reflect.ValueOf(req.Param).Elem())
	//kind := paramType.Kind()
	//fmt.Println(kind)
	//fmt.Println(paramType)
	//fmt.Println(paramType)
	switch paramType.Kind() {
	case reflect.Struct:
		//fmt.Println("num of fields: ", param.NumField())
		for i := 0; i < param.NumField(); i++ {
			//fmt.Println(param.Field(i).Kind())
			//fmt.Println(param.Field(i))
			//fmt.Println()
			fieldType := param.Field(i).Kind()
			switch fieldType {
			case reflect.Int32:
				buf[offset] = 8
				offset += 1
				break
			}
			offset += Int16ToBytes(int16(i+1), buf[offset:])
			offset += Int32ToBytes(param.Field(i).Interface().(int32), buf[offset:])
		}
		break
	}
	//enc.w.Write(Int32ToBytes(version))
	buf[offset] = 0
	offset += 1

	_, err := enc.w.Write(buf[0: offset])
	if err != nil {
		fmt.Println(err.Error())
	}
	//fmt.Println("sent ", n, " bytes")
	return nil
}