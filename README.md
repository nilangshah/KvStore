KvStore
=======

It's a key value store written in GO Language. 

Usage:
    
    go get github.com/nilangshah/KvStore
    
    go install github.com/nilangshah/KvStore
    
    bin/janGoDB
    
    This will start a server which is listening on port 11211 on localhost.   
    
    Now go in test folder which contains a java file which is created for testing of this keyvalue store. 
    
    javac tcpclient.java
    java TcpClient
    go run test1_go.go
        This will start 100 client and will set and get values concurrently.
    go run test2_go.go
        This is a normal test for set and get method.
        
    You will see print statements of sending and recieving massages between client and server.
    
    I have assumed that you have set GOPATH and PATH for GOLANG.
    Currently I have written client only in java but I will add client in Go and python for testing purpose.
    
    Thank you. 
