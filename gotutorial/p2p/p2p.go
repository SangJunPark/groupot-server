package p2p

import (
	"fmt"
	"gotutorial/blockchain"
	"gotutorial/utils"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader
var conns []*websocket.Conn

func Upgrade(rw http.ResponseWriter, r *http.Request) {

	ip := utils.Splitter(r.RemoteAddr, ":", 0)
	openPort := r.URL.Query().Get("openPort")

	fmt.Printf("Upgrade address ip %s openport %s\n", ip, openPort)

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return ip != "" || openPort != ""
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	conns = append(conns, conn)
	utils.HandleErr(err)

	// d[1] - is not an open port
	initPeer(conn, ip, openPort)

	// for {
	// 	_, p, err := conn.ReadMessage()
	// 	if err != nil {
	// 		//conn.Close()
	// 		break
	// 	}
	// 	for _, aConn := range conns {
	// 		if conn != aConn {
	// 			aConn.WriteMessage(websocket.TextMessage, p)
	// 		}
	// 	}
	// }
}

func AddPeer(address, port, openPort string, broadcast bool) {

	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil)
	utils.HandleErr(err)
	peer := initPeer(conn, address, port)

	fmt.Printf("AddPeer %s   %s   %s", address, port, openPort)

	if broadcast {
		BroadcastNewPeer(peer)
		return
	}
	sendNewestBlock(peer)
}

func BroadcastNewBlock(block *blockchain.Block) {
	for _, p := range Peers.v {
		notifyNewBlock(p, block)
	}
}

func BroadcastNewTransaction(tx *blockchain.Tx) {
	for _, p := range Peers.v {
		notifyNewTransaction(p, tx)
	}
}

func BroadcastNewPeer(newP *peer) {
	for _, p := range Peers.v {
		if newP.id != p.id {
			payload := fmt.Sprintf("%s:%s", newP.id, p.port)
			notifyNewPeer(p, payload)
		}
	}
}
