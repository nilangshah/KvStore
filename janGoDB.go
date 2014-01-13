package main

import (
        "bufio"
        "flag"
        "fmt"
        "io"
        "net"
        "strings"
)

var (
  kvMap = map[string][]byte{}
)

func main() {
            flag.Parse()

        listener, err := net.Listen("tcp", "127.0.0.1:11211")
        if err != nil {
                panic("Error listening on 11211: " + err.Error())
        }

        kvMap = make(map[string][]byte)

                for {
                        netconn, err := listener.Accept()
                        if err != nil {
                                panic("Accept error: " + err.Error())
                        }

                        go handleConn(netconn)
                }

}

/*
* Networking
*/
func handleConn(conn net.Conn) {
    defer conn.Close()
        reader := bufio.NewReader(conn)
        for {

                // Fetch

                content, err := reader.ReadString('\n')
                if err == io.EOF {
                        break
                } else if err != nil {
                        fmt.Println(err)
                        return
                }

                content = content[:len(content)-1] // Chop \n

                // Handle

                subContent := strings.Split(content, " ")
                cmd := subContent[0]
                switch cmd {

                case "get":
                        key := subContent[1]
                        val, ok := kvMap[key]
                        fmt.Println("value get = ",string(val))
                        if ok {
                                conn.Write([]uint8(string(val) + "\r\n"))
                        }
                      continue
                case "set":
                        key := subContent[1]
                        fmt.Println(key)
                        val := subContent[2]
                        kvMap[key] = []byte(val)
                        conn.Write([]uint8("STORED\r\n"))
			continue
		case "quit":
			break

                }
	break
        }
}
