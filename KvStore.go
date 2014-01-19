package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
)

type Configurations struct {
	Path     string
	FilePerm os.FileMode
	PathPerm os.FileMode
}

type InDex struct {
	off  int64
	size int64
}

var inDexFileName string = "inDex.0"
var dbFileName string = "KvStore.0"

type GoDB struct {
	sync.RWMutex
	Configurations
	kvMap     map[string][]byte
	inDexMap  map[string][]int64
	dbFile    *os.File
	inDexFile *os.File
}

const (
	defaultPath                 = "KvStore"
	defaultFilePerm os.FileMode = 0666
	defaultPathPerm os.FileMode = 0777
)

var goDB *GoDB

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:11211")
	if err != nil {
		panic("Error listening on 11211: " + err.Error())
	}

	goDB = New(Configurations{
		Path:     "my-KvStore",
		FilePerm: defaultFilePerm,
		PathPerm: defaultPathPerm,
	})
	for {
		netconn, err := listener.Accept()
		if err != nil {
			panic("Accept error: " + err.Error())
		}

		go handleConn(netconn)
	}

}

func New(config Configurations) *GoDB {
	d := &GoDB{
		Configurations: config,
		kvMap:          map[string][]byte{},
	}
	return d
}

func OpenDBFiles() {
	goDB.Lock()
	defer goDB.Unlock()
	var err error
	filePath := fmt.Sprintf("%s%c%s", goDB.Path, os.PathSeparator, dbFileName)
	indexPath := fmt.Sprintf("%s%c%s", goDB.Path, os.PathSeparator, inDexFileName)
	mode := os.O_RDWR | os.O_CREATE

	if err = os.MkdirAll(goDB.Path, goDB.PathPerm); err != nil {
		panic("CREATE DIR FAILED:" + err.Error())
		return
	}
	goDB.dbFile, err = os.OpenFile(filePath, mode, goDB.FilePerm)
	if err != nil {
		panic("OPEN FILE FAILED:" + err.Error())
		return
	}

	goDB.inDexFile, err = os.OpenFile(indexPath, mode, goDB.FilePerm)
	if err != nil {
		panic("OPEN FILE FAILED:" + err.Error())
		return
	}
	state, _ := os.Stat(getFilePath(inDexFileName))
	if state.Size() > 0 && goDB.inDexMap == nil {
		goDB.inDexMap = make(map[string][]int64)
		dec := gob.NewDecoder(goDB.inDexFile)
		err = dec.Decode(&goDB.inDexMap)
		fmt.Println(len(goDB.inDexMap))
		fmt.Println(len(goDB.kvMap))
		if err != nil {
			panic("DECODE ERRROR: " + err.Error())
			return
		}
	}
}

/*
* Networking
 */
func handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	OpenDBFiles()

	//goDB.dbFile.Close()
	//goDB.inDexFile.Close()
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
				//fmt.Println("value found in main mem")
				goDB.RUnlock()
				conn.Write([]uint8(string(val) + "\r\n"))

			} else {
				//fmt.Println("not found in main mem")
				val1, err := ReadFromDisk(key)
				if err != nil {
					goDB.RUnlock()
					conn.Write([]uint8("ERROR\r\n"))
				} else {

					goDB.RUnlock()
					goDB.Lock()
					goDB.kvMap[key] = []byte(val1)
					goDB.Unlock()
					conn.Write([]uint8(string(val1) + "\r\n"))
				}
			}
			continue
		case "set":
			goDB.Lock()
			key := subContent[1]
			val := subContent[2]

			err := WriteToDisk(1, key, val)

			if err != nil {
				panic(err.Error())
				conn.Write([]uint8("ERROR OCCURED\r\n"))
			} else {
				goDB.kvMap[key] = []byte(val)
				conn.Write([]uint8("STORED\r\n"))
			}
			goDB.Unlock()
			continue
		case "delete":
			goDB.Lock()
			key := subContent[1]
			err := WriteToDisk(-1, key, "\n")
			if err != nil {
				panic(err.Error())
				conn.Write([]uint8("ERROR OCCURED\r\n"))
			} else {
				delete(goDB.kvMap, key)
				conn.Write([]uint8("DELETED\r\n"))
			}
			goDB.Unlock()
			continue
		case "quit":
			conn.Write([]uint8("Bye\r\n"))
			break

		}
		break
	}
}

func getFilePath(name string) string {
	filePath := fmt.Sprintf("%s%c%s", goDB.Path, os.PathSeparator, name)
	return filePath
}

func WriteToDisk(op int, key string, val string) error {

	filePath := getFilePath(dbFileName)
	fInfo, _ := os.Stat(filePath)
	//fmt.Println(fInfo.Size())
	if op == 1 {
		a := string(len(key)) + string(len(val)) + string(key) + string(val)
		n1 := fInfo.Size()
		n, err := goDB.dbFile.WriteAt([]byte(a), n1)
		arr := make([]int64, 2)
		arr[0] = int64(n)  //size
		arr[1] = int64(n1) //offset
		if goDB.inDexMap == nil {
			goDB.inDexMap = make(map[string][]int64)
		}
		goDB.inDexMap[key] = arr
		//fmt.Println("index map"+[]byte(len(goDB.inDexMap)))
		//fmt.Println("kvmap "+[]byte(len(goDB.kvMap)))
		b := new(bytes.Buffer)
		enc := gob.NewEncoder(b)
		err = enc.Encode(goDB.inDexMap)
		if err != nil {
			return err
		}
		e := ioutil.WriteFile(getFilePath(inDexFileName), b.Bytes(), 0666)
		if e != nil {
			return e
		}
		if err != nil {
			return err
		}
	} else {
		a := "0" + string(len(key)) + string(key)
		n1 := fInfo.Size()
		_, err := goDB.dbFile.WriteAt([]byte(a), n1)
		if err != nil {
			return err
		} else {
			delete(goDB.inDexMap, key)
			b := new(bytes.Buffer)
			enc := gob.NewEncoder(b)
			err = enc.Encode(goDB.inDexMap)
			if err != nil {
				return err
			}
			e := ioutil.WriteFile(getFilePath(inDexFileName), b.Bytes(), 0666)
			if e != nil {
				return e
			}

			return nil
		}
	}
	return nil

}

func ReadFromDisk(key string) (string, error) {
	//fmt.Println(fInfo.Size())

	if n, ok := goDB.inDexMap[key]; ok {
		//	fmt.Println("not found in index map")
		dat := make([]byte, n[0])
		_, err := goDB.dbFile.ReadAt(dat, n[1])

		if err != nil {
			return "NIL", err
		}

		return string(dat[(2 + dat[0]):(2 + dat[0] + dat[1])]), nil
	} else {
		//	fmt.Println("not found in index map")
		return "NIL", nil //value not exist
	}
	return "NIL", nil

}
