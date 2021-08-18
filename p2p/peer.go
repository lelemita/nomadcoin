package p2p

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var Peers map[string]*peer = make(map[string]*peer)

type peer struct {
	key     string
	address string
	port    string
	conn    *websocket.Conn
	inbox   chan []byte
}

func (p *peer) close() {
	p.conn.Close()
	delete(Peers, p.key)
}

func (p *peer) read() {
	defer p.close()
	for {
		_, m, err := p.conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("%s", m)
	}
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
}

func initPeer(conn *websocket.Conn, address, port string) *peer {
	key := fmt.Sprintf("%s:%s", address, port)
	p := &peer{
		key:     key,
		address: address,
		port:    port,
		conn:    conn,
		inbox:   make(chan []byte),
	}
	Peers[key] = p
	go p.read()
	go p.write()
	return p
}
