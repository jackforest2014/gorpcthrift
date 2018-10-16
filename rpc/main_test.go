package main

import (
	"context"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"sync"
	"techValidate/thriftRPC"
	"techValidate/tutorial"
	"testing"
	"time"
)

//var conn net.Conn = nil

func Benchmark_RpcThriftTest(t *testing.B){
	rand.Seed(time.Now().Unix())


	conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		log.Fatal("net.Dial:", err)
	}

	fmt.Println(conn.LocalAddr())
	client := rpc.NewClientWithCodec(thriftRPC.NewThriftCodec(conn))

	//for i := 1; i <= 1; i++ {
		makeRPCThriftCall(client)
	//}

	conn.Close()
}

func makeRPCThriftCall(client *rpc.Client){


	wg := sync.WaitGroup{}

	//a := int32(100000) //rand.Int31() % 1000
	//b := int32(10000) //rand.Int31() % 1000

	//for i := 1; i <= 33750; i++ {
	for i := 1; i <= 1000; i++ {
		//fmt.Println("# ", i, " --------------------", client)
		go func(i int) {
			wg.Add(1)
			a := rand.Int31() % 1000
			b := rand.Int31() % 1000
			reply := result{}
			//fmt.Println("#", i, "   a = ", a, "  b = ", b)
			err := client.Call("add", &request{Num1: a, Num2: b}, &reply)
			if err != nil {
				log.Fatal(err)
			}
			//if a + b != reply.Result {
			//	log.Fatal(a, " + ", b, " != ", reply.Result, " !")
			//}
			fmt.Println("#", i, "   ", a, " + ", b, " = ", reply.Result, "  --> ", a + b == reply.Result)
			//fmt.Println()
			//fmt.Println(reply.Result)
			wg.Done()
		}(i)
		//time.Sleep(time.Millisecond * 100)
	}
	//time.Sleep(time.Second*2)
	wg.Wait()
}

func Benchmark_BothThriftCall(b *testing.B) {
	makeThriftCall()
}

func makeThriftCall() {
	socket, err := thrift.NewTSocketTimeout(":9090", time.Second)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("%v\n", socket.Addr())
	trans := thrift.NewTFramedTransport(socket)
	//defer trans.Close()

	protocolFactory := thrift.NewTCompactProtocolFactory()
	inProto := protocolFactory.GetProtocol(trans)
	outProto := protocolFactory.GetProtocol(trans)

	client := tutorial.NewCalculatorClient(thrift.NewTStandardClient(inProto, outProto))

	invokeCtx, f := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*10000))
	defer f()

	a := rand.Int31() % 1000
	b := rand.Int31() % 1000
	result, err := client.Add(invokeCtx, a, b)
	if a + b != result {
		log.Fatal(a, " + ", b, " != ", result)
	}

	//socket.Close()

}