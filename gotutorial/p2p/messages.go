package p2p

import (
	"encoding/json"
	"fmt"
	"gotutorial/blockchain"
	"gotutorial/utils"
)

type MessageKind int

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlockResponse
	MessageNewBlockCreated
	MessageNewTransactionCreated
	MessageNewPeerConnected
)

type Message struct {
	Kind    MessageKind
	Payload []byte
}

func sendNewestBlock(p *peer) {
	fmt.Println("sendnewestblock")
	b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
	utils.HandleErr(err)

	m := makeMessage(MessageNewestBlock, b)
	p.inbox <- m
}

func sendAllBlocks(p *peer) {
	fmt.Println("sendallblocks")
	b := blockchain.Blocks(blockchain.Blockchain())
	m := makeMessage(MessageAllBlockResponse, b)
	p.inbox <- m
}

func requestAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksRequest, nil)
	p.inbox <- m
}

func notifyNewBlock(p *peer, block *blockchain.Block) {
	m := makeMessage(MessageNewBlockCreated, block)
	p.inbox <- m
}

func notifyNewTransaction(p *peer, tx *blockchain.Tx) {
	m := makeMessage(MessageNewTransactionCreated, tx)
	p.inbox <- m
}

func notifyNewPeer(p *peer, address string) {
	m := makeMessage(MessageNewPeerConnected, address)
	p.inbox <- m
}

func (m *Message) addPayload(p interface{}) {
	b, err := json.Marshal(p)
	utils.HandleErr(err)
	m.Payload = b
}

func makeMessage(kind MessageKind, payload interface{}) []byte {
	m := Message{
		Kind:    kind,
		Payload: utils.ToJSON(payload),
	}
	return utils.ToJSON(m)
}

func handleMessage(m *Message, p *peer) {
	fmt.Printf("%d %s\n", m.Kind, p.id)

	switch m.Kind {
	case MessageNewestBlock:
		var payload blockchain.Block
		json.Unmarshal(m.Payload, &payload)
		fmt.Println(payload.Hash)
		b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
		utils.HandleErr(err)
		if payload.Height <= b.Height {
			sendAllBlocks(p)
		} else {
			sendNewestBlock(p)
		}
	case MessageAllBlocksRequest:
		sendAllBlocks(p)
	case MessageAllBlockResponse:
		var payload []*blockchain.Block
		json.Unmarshal(m.Payload, &payload)
		blockchain.Blockchain().Replace(payload)
		fmt.Println("received allblock response")
	case MessageNewBlockCreated:
		var payload *blockchain.Block
		json.Unmarshal(m.Payload, &payload)
		blockchain.Blockchain().AddPeerBlock(payload)
	case MessageNewTransactionCreated:
		var payload *blockchain.Tx
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Mempool().AddPeerTx(payload)
		fmt.Println(payload)
	case MessageNewPeerConnected:
		var payload string
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		fmt.Println("MessageNewPeerConnected " + payload)
	default:
		fmt.Println("Undefined message kind")
	}
}
