package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:11211")
	defer conn.Close()
	if err != nil {
		panic("Connection error: " + err.Error())
	}
	reader := bufio.NewReader(conn)
	sentence := "set alpha gamma\n"
	fmt.Println("SEND SERVER:" + sentence)
	conn.Write([]uint8(sentence))
	content, err := reader.ReadString('\n')
	if err == io.EOF {
	} else if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("FROM SERVER:" + content)
	sentence = "set alpha1 gamma1\n"

	fmt.Println("SEND SERVER:" + sentence)
	conn.Write([]uint8(sentence))
	content, err = reader.ReadString('\n')
	if err == io.EOF {
	} else if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("FROM SERVER:" + content)
	sentence = "delete alpha\n"

	fmt.Println("SEND SERVER:" + sentence)
	conn.Write([]uint8(sentence))
	content4, err := reader.ReadString('\n')
	if err == io.EOF {
	} else if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("FROM SERVER:" + content4)

	sentence = "get alpha\n"

	fmt.Println("SEND SERVER:" + sentence)
	conn.Write([]uint8(sentence))
	content1, err := reader.ReadString('\n')
	if err == io.EOF {
	} else if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("FROM SERVER:" + content1)

	sentence = "get alpha1\n"

	fmt.Println("SEND SERVER:" + sentence)
	conn.Write([]uint8(sentence))
	content2, err := reader.ReadString('\n')
	if err == io.EOF {
	} else if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("FROM SERVER:" + content2)
	sentence = "quit" + "\n"
	fmt.Println("SEND SERVER:" + sentence)
	conn.Write([]uint8(sentence))
	content5, err := reader.ReadString('\n')
	if err == io.EOF {
	} else if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("FROM SERVER:" + content5)

	fmt.Println("Thank you Test successfull")
}
