package main

import (
        "bufio"
        "fmt"
        "io"
        "net"
	"sync"
        "strings"
)



type GoDB struct {
	sync.RWMutex
	kvMap     map[string][]byte
}
 
var goDB *GoDB

func main() {
        listener, err := net.Listen("tcp", "127.0.0.1:11211")
        if err != nil {
                panic("Error listening on 11211: " + err.Error())
        }

        goDB = New()
         for {
                        netconn, err := listener.Accept()
                        if err != nil {
                                panic("Accept error: " + err.Error())
                        }

                        go handleConn(netconn)
             }

}


func New() *GoDB {
	d := &GoDB{
		kvMap:     map[string][]byte{},
	}

	return d
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
			goDB.RLock()
			key := subContent[1]
                        val, ok := goDB.kvMap[key]
                        if ok {
                                conn.Write([]uint8(string(val) + "\r\n"))
                        } else {
				 conn.Write([]uint8("NIL\r\n"))		
			}
			goDB.RUnlock()
                      continue
                case "set":
			goDB.Lock()
                        key := subContent[1]
                        val := subContent[2]
                        goDB.kvMap[key] = []byte(val)
                        conn.Write([]uint8("STORED\r\n"))
			goDB.Unlock()
			continue
		case "quit":
			break

                }
	break
        }
}
