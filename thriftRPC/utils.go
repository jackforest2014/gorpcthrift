package thriftRPC

import "encoding/binary"

func Int16ToBytes(i int16, buf []byte) int {
	binary.BigEndian.PutUint16(buf, uint16(i))
	//buf[0] = int8(255 & (i >> 8))
	//buf[1] = int8(255 & i)
	return 2
}

func BytesToInt16(buf[] byte) int16{
	return int16(binary.BigEndian.Uint16(buf))
}


func Int32ToBytes(i int32, buf []byte) int {
	//buf[0] = int8(255 & (i >> 24))
	//buf[1] = int8(255 & (i >> 16))
	//buf[2] = int8(255 & (i >> 8))
	//buf[3] = int8(255 & i)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return 4
}

func BytesToInt32(buf[] byte) int32 {
	//return int32((uint32(buf[0]) << 24) + (uint32(buf[1]) << 16) + (uint32(buf[2]) << 8) + uint32(buf[3]))
	return int32(binary.BigEndian.Uint32(buf))
}

func UInt64ToBytes(i uint64, buf []byte) int {
	binary.BigEndian.PutUint64(buf, i)
	//buf[0] = int8(255 & (i >> 56))
	//buf[1] = int8(255 & (i >> 48))
	//buf[2] = int8(255 & (i >> 40))
	//buf[3] = int8(255 & (i >> 32))
	//buf[4] = int8(255 & (i >> 24))
	//buf[5] = int8(255 & (i >> 16))
	//buf[6] = int8(255 & (i >> 8))
	//buf[7] = int8(255 & i)
	return 8
}

func StringToBytes(s string, buf []byte) int {
	strBytes := []byte(s)
	for i, b := range strBytes {
		buf[i] = b
	}
	return len(strBytes)
}