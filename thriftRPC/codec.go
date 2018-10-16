package thriftRPC

import (
	"errors"
	"fmt"
	"io"
	"net/rpc"
	"sync"
)

type ThriftCodec struct {
	dec *ThriftDecoder // for reading JSON values
	enc *ThriftEncoder // for writing JSON values
	c   io.Closer

	//// temporary work space
	req  thriftRequest
	//resp clientResponse

	// JSON-RPC responses include the request id but not the request method.
	// Package rpc expects both.
	// We save the request method in pending when sending a request
	// and then look it up by request ID when filling out the rpc Response.
	mutex   sync.Mutex        // protects pending
	pending map[uint64]string // map request id to method name

	//Rwc		io.ReadWriteCloser
}

type thriftRequest struct {
	Method string         	`json:"method"`
	Param  interface{} 		`json:"param"`
	Id     uint64         	`json:"id"`
}


func NewThriftCodec(conn io.ReadWriteCloser) rpc.ClientCodec {
	return &ThriftCodec{
		dec:     NewThriftDecoder(conn),
		enc:     NewThriftEncoder(conn),
		c:       conn,
		pending: make(map[uint64]string),
	}
}

func (c *ThriftCodec) WriteRequest(r *rpc.Request, body interface{}) (err error){
	if body == nil {
		return fmt.Errorf("Nil request body from client.")
	}

	//fmt.Println("seq: ", r.Seq)
	c.pending[r.Seq] = r.ServiceMethod
	c.req.Id = r.Seq
	r.Seq += 2
	c.req.Method = r.ServiceMethod
	c.req.Param = body

	//query := body.(*Query)
	return c.enc.Encode(&c.req)
}


func (c *ThriftCodec) ReadResponseHeader(r *rpc.Response) error {
	version := readInt32(c.dec.r)
	var methodName string
	//var typ uint8 = uint8(version & 255)
	var seqId int32

	if version < 0 {
		version = version & -65536
		if version != -2147418112 {
			return errors.New("bad version in response message")
		} else {
			methodName = readString(c.dec.r)
			seqId = readInt32(c.dec.r)
		}
	}
	//fmt.Println("version = ", version, ", method = ", methodName, ", type = ", typ, ", seqID = ", seqId)
	r.Seq = uint64(seqId)
	r.ServiceMethod = methodName

	return nil
}

func (c *ThriftCodec) ReadResponseBody(x interface{}) error {
	if x == nil {
		return nil //errors.New("invalid response")
	}

	c.dec.Decode(x)
	return nil
}

func (c *ThriftCodec) Close() error {
	return nil
}