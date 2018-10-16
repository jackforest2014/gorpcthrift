package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	//simple()
	json()
}

func simple(){
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var reply string
	err = client.Call("HelloService.Hello", "你好", &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}

func json(){
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("net.Dial:", err)
	}
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	var reply string
	err = client.Call("HelloService.Hello", "你好", &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}

type result struct{
	Result int32
}

type request struct {
	Num1  int32
	Num2  int32
}

func json2Java(){
	conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		log.Fatal("net.Dial:", err)
	}

	//for i := 1; i <= 20; i++ {
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	var reply string

	//go func() {
	a := rand.Int31() % 1000
	b := rand.Int31() % 1000
	fmt.Println("a = ", a, "  b = ", b)
	err = client.Call("add", &request{Num1: a, Num2: b}, &result{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
	//}()
	//}
}
