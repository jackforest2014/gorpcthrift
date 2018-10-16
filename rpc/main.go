package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"reflect"
	"sync"
	"techValidate/thriftRPC"
	"time"
)

type result struct{
	Result int32
}

type request struct {
	Num1  int32
	Num2  int32
}

var clientPool sync.Pool

func main() {
	fmt.Println("RPC")

	req := request{Num1:1,Num2:2}
	var reqInterface interface{}
	reqInterface = req
	typ := reflect.TypeOf(reqInterface)
	val := reflect.ValueOf(reqInterface)
	fmt.Println(typ.Name())
	fmt.Println(typ.NumField())
	fmt.Println(typ.NumMethod())

	fmt.Println(typ.Kind())
	for i:=0; i < typ.NumField(); i++{
		fmt.Print(typ.Field(i).Type.Name(), " - ", typ.Field(i).Type.Kind(), " : ", val.Field(i))
		fmt.Println()
	}

	fmt.Println(typ)
	fmt.Println(val)

	rand.Seed(time.Now().Unix())

	for i := 1; i <= 1; i++ {
		thriftCall()
	}

	//client, err := rpc.Dial("tcp", ":9090")
	//if err != nil {
	//	log.Fatal("dial: ", err.Error())
	//}
	//
	//err2 := client.Call("add", &request{Num1:3, Num2:4}, &result{})
	//if err2 != nil {
	//	log.Fatal("call ping: ", err2.Error())
	//}
}


func thriftCall(){
	conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		log.Fatal("net.Dial:", err)
	}

	fmt.Println(conn.LocalAddr())
	client := rpc.NewClientWithCodec(thriftRPC.NewThriftCodec(conn))

	wg := sync.WaitGroup{}

	for i := 1; i <= 1000; i++ {
		//fmt.Println("# ", i, " --------------------", client)
		go func(i int) {
			wg.Add(1)
			a := rand.Int31() % 1000
			b := rand.Int31() % 1000
			reply := result{}
			fmt.Println("#", i, "   a = ", a, "  b = ", b)
			err = client.Call("add", &request{Num1: a, Num2: b}, &reply)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("#", i, "   ", a, " + ", b, " = ", reply.Result, "  --> ", a + b == reply.Result)
			//fmt.Println()
			//fmt.Println(reply.Result)
			wg.Done()
		}(i)
		//time.Sleep(time.Millisecond * 100)
	}
	//time.Sleep(time.Second*2)
	wg.Wait()
	conn.Close()
}