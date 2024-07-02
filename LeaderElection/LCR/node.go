package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"strconv"
	"time"
)

var nodeAddressMap = map[int]string{
	0: "localhost:5000",
	1: "localhost:5001",
	2: "localhost:5002",
	3: "localhost:5003",
	4: "localhost:5004",
}

type Node struct {
	ID      int
	Address string
	Peers   *Peers
	UUID    int
	Status  string
}

func NewNode(nodeId int) *Node {
	node := &Node{
		ID:      nodeId,
		Address: nodeAddressMap[nodeId],
		Peers:   NewPeers(),
		UUID:    rand.Intn(1000),
		Status:  "COMMON",
	}
	return node
}

func (node *Node) NewListener() (net.Listener, error) {
	addr, err := net.Listen("tcp", node.Address)
	return addr, err
}

func (node *Node) ConnectToPeers() {
	fmt.Println("My UUID is ", node.UUID)
	for nodeId, nodeAddress := range nodeAddressMap {
		if node.ID == nodeId {
			continue
		}

		rpcClient := node.connect(nodeAddress)
		pingMessage := Message{
			FromPeerId: node.ID,
			Type:       "PING",
			Payload:    "TESTPING",
		}

		reply, _ := node.SendMessage(rpcClient, pingMessage)

		fmt.Println("received message from ", reply.FromPeerId, " of type ", reply.Type, " containing payload ", reply.Payload)
		if reply.Type == "PONG" {
			node.Peers.Add(nodeId, rpcClient)
		}
	}
}

func (node *Node) connect(peerAddress string) *rpc.Client {
retry:
	client, err := rpc.Dial("tcp", peerAddress)
	if err != nil {
		fmt.Println("Error dialing rpc ", err.Error())
		time.Sleep(100 * time.Millisecond)
		goto retry
	}
	return client
}

func (node *Node) SendMessage(RPCClient *rpc.Client, args Message) (Message, error) {
	var reply Message
	err := RPCClient.Call("Node.HandleMessage", args, &reply)
	if err != nil {
		fmt.Println("Error sending message ", err.Error())
	}
	return reply, err
}

func (node *Node) HandleMessage(args Message, reply *Message) error {
	fmt.Println("received message from ", args.FromPeerId, " of type ", args.Type, " containing payload ", args.Payload)
	reply.FromPeerId = node.ID
	switch args.Type {
	case "PING":
		reply.Type = "PONG"
		if args.Payload == "TESTPING" {
			reply.Payload = "TESTPONG"
		} else if args.Payload == "LEADERPING" {
			reply.Payload = "LEADERPONG"
		}

	case "ELECTION":
		messageUUID, _ := strconv.Atoi(args.Payload)
		if messageUUID > node.UUID { // need to propagate message
			neighbours := node.getNeighbours()
			// send election message with payload to neighbour (unidirectional)
			for _, neighbour := range neighbours {
				peer := node.Peers.Get(neighbour)
				node.SendMessage(peer.RPCClient, args)
			}
		} else if messageUUID == node.UUID { // this node is the leader
			node.Status = "LEADER"
		}
	}

	return nil
}

func (node *Node) elect() {
	neighbours := node.getNeighbours()

	// send election message with payload to neighbour (unidirectional)
	for _, neighbour := range neighbours {
		electionMessage := Message{
			FromPeerId: node.ID,
			Type:       "ELECTION",
			Payload:    strconv.Itoa(node.UUID),
		}
		peer := node.Peers.Get(neighbour)
		node.SendMessage(peer.RPCClient, electionMessage)
	}

	// Sleep and respond to messages periodically.
idle:
	if node.Status == "COMMON" {
		fmt.Println("I have woken up to handle messages, will go back to sleep after processing them if any")
	} else if node.Status == "LEADER" {
		fmt.Println("I am the leader, I will send PING to all")

		peers := node.Peers.GetPeers()
		for _, peer := range peers {
			pingMessage := Message{
				FromPeerId: node.ID,
				Type:       "PING",
				Payload:    "LEADERPING",
			}
			reply, _ := node.SendMessage(peer.RPCClient, pingMessage)
			fmt.Println("received message from ", reply.FromPeerId, " of type ", reply.Type, " containing payload ", reply.Payload)
		}
	}
	time.Sleep(3 * time.Second)
	goto idle

}

// LCR is a unidirectional ring. So each node connects to its immediate neighbour in ACW direction only (during election)
func (node *Node) getNeighbours() []int {
	neighbours := make([]int, 0)
	neighbours = append(neighbours, (node.ID+1)%len(nodeAddressMap))
	return neighbours
}
