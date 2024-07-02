package main

import (
	"fmt"
	"net/rpc"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Invalid args, enter some id")
		return
	}
	nodeId, _ := strconv.Atoi(os.Args[1])

	node := NewNode(nodeId)
	listener, err := node.NewListener()
	if err != nil {
		fmt.Println("Error listening %s", err.Error())
		return
	}

	defer listener.Close()

	rpcServer := rpc.NewServer()
	rpcServer.Register(node)

	go rpcServer.Accept(listener)

	node.ConnectToPeers()

	// Wait for some time before triggering election
	time.Sleep(15 * time.Second)

	node.elect()
}
