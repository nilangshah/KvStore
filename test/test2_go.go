package main

import (
    "fmt"
	"net"
	"bufio"
	"io"
)

func main() {
        conn, err := net.Dial("tcp", "127.0.0.1:11211")
	if err != nil {
	 panic("Connection error: " + err.Error())
	}
	reader := bufio.NewReader(conn)
	sentence := "set alpha beta\n"
	fmt.Println("SEND SERVER:" + sentence)	
	conn.Write([]uint8(sentence))
	content, err := reader.ReadString('\n')
                if err == io.EOF {
                } else if err != nil {
                        fmt.Println(err)
                        return
                }
	fmt.Println("FROM SERVER:"+content)
	sentence = "get alpha\n"
	
	fmt.Println("SEND SERVER:" + sentence)	
	conn.Write([]uint8(sentence))
	content, err = reader.ReadString('\n')
                if err == io.EOF {
                } else if err != nil {
                        fmt.Println(err)
                        return
                }
	fmt.Println("FROM SERVER:"+content)

	sentence = "get gamma\n"
	
	fmt.Println("SEND SERVER:" + sentence)	
	conn.Write([]uint8(sentence))
	content, err = reader.ReadString('\n')
                if err == io.EOF {
                } else if err != nil {
                        fmt.Println(err)
                        return
                }
	fmt.Println("FROM SERVER:"+content)



	fmt.Println("Thank you Test successfull")
}
