package main

import "net/rpc"
import "sync"

type Peer struct {
	ID        int
	RPCClient *rpc.Client
}

type Peers struct {
	*sync.RWMutex
	peerClientMap map[int]*Peer
}

func NewPeers() *Peers {
	return &Peers{
		RWMutex:       &sync.RWMutex{},
		peerClientMap: make(map[int]*Peer),
	}
}

func (p *Peers) Add(ID int, client *rpc.Client) {
	p.Lock()
	defer p.Unlock()

	p.peerClientMap[ID] = &Peer{ID: ID, RPCClient: client}
}

func (p *Peers) Get(ID int) *Peer {
	p.Lock()
	defer p.Unlock()

	val := p.peerClientMap[ID]
	return val
}

func (p *Peers) GetPeers() []Peer {
	p.Lock()
	defer p.Unlock()

	peers := make([]Peer, 0, len(p.peerClientMap))
	for _, peer := range p.peerClientMap {
		peers = append(peers, *peer)
	}

	return peers
}

func (p *Peers) GetPeerIds() []int {
	p.Lock()
	defer p.Unlock()

	peerIDs := make([]int, 0, len(p.peerClientMap))
	for _, peer := range p.peerClientMap {
		peerIDs = append(peerIDs, peer.ID)
	}

	return peerIDs
}
