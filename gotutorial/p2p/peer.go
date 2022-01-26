package p2p

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

var Peers peers = peers{
	v: make(map[string]*peer),
}

type peer struct {
	id     string
	conn   *websocket.Conn
	inbox  chan []byte
	closed bool
}

type peers struct {
	v map[string]*peer
	m sync.Mutex
}

func AllPeers(p *peers) []string {
	p.m.Lock()
	defer p.m.Unlock()

	var peers []string
	for key := range p.v {
		peers = append(peers, key)
	}

	return peers
}

func initPeer(conn *websocket.Conn, address, port string) *peer {
	// Peers.m.Lock()
	// defer Peers.m.Unlock()
	key := fmt.Sprintf("%s:%s", address, port)

	p := &peer{
		key,
		conn,
		make(chan []byte),
		false,
	}
	Peers.v[key] = p
	go p.read()
	go p.write()
	return p
}

func (p *peer) read() {
	defer p.close()
	for {
		m := Message{}
		err := p.conn.ReadJSON(&m)
		if err != nil {
			fmt.Println(err)
			break
		}

		handleMessage(&m, p)
	}
	fmt.Println("read done - peer closed")
}

func (p *peer) write() {
	defer p.close()
	for {
		m, ok := <-p.inbox
		if !ok {
			break
		}
		p.conn.WriteMessage(websocket.TextMessage, m)
	}
	fmt.Println("write done - peer closed")
}

func (p *peer) close() {
	defer Peers.m.Unlock()

	Peers.m.Lock()
	if p.closed {
		return
	}

	close(p.inbox)
	p.conn.Close()
	delete(Peers.v, p.id)
	p.closed = true
}
