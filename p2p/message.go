package p2p

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lelemita/nomadcoin/blockchain"
	"github.com/lelemita/nomadcoin/utils"
)

type MessageKind int

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResponse
	MessageNewBlockNotify
	MessageNewTxNotify
	MessageNewPeerNotify
)

type Message struct {
	Kind    MessageKind
	Payload []byte
}

func makeMessage(kind MessageKind, payload interface{}) []byte {
	m := Message{
		Kind:    kind,
		Payload: utils.ToJson(payload),
	}
	return utils.ToJson(m)
}

func sendNewestBlock(p *peer) {
	fmt.Printf("Sending newest block to %s\n", p.key)
	b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
	utils.HandleErr(err)
	m := makeMessage(MessageNewestBlock, b)
	p.inbox <- m
}

func requestAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksRequest, nil)
	p.inbox <- m
}

func sendAllBlocks(p *peer) {
	blocks := blockchain.Blocks(blockchain.Blockchain())
	p.inbox <- makeMessage(MessageAllBlocksResponse, blocks)
}

func notifyNewBlock(p *peer, b *blockchain.Block) {
	p.inbox <- makeMessage(MessageNewBlockNotify, b)
}

func notifyNewTx(p *peer, tx *blockchain.Tx) {
	fmt.Printf("Notify New Tx to %s\n", p.port)
	p.inbox <- makeMessage(MessageNewTxNotify, tx)
}

func notifyNewPeer(p *peer, address string) {
	fmt.Printf("Notify New Peer to %s\n", p.port)
	p.inbox <- makeMessage(MessageNewPeerNotify, address)
}

func handleMsg(m *Message, p *peer) {
	switch m.Kind {
	case MessageNewestBlock:
		fmt.Printf("Received the newest block from %s\n", p.key)
		var payload blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
		utils.HandleErr(err)
		if payload.Height >= b.Height {
			fmt.Printf("Requesting all blocks from %s\n", p.key)
			requestAllBlocks(p)
		} else {
			fmt.Printf("I'll Send newest block to %s\n", p.key)
			sendNewestBlock(p)
		}
	case MessageAllBlocksRequest:
		fmt.Printf("%s wants all the blocks.\n", p.key)
		sendAllBlocks(p)
	case MessageAllBlocksResponse:
		fmt.Printf("Received all the blocks from %s\n", p.key)
		var payload []*blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		// todo: ?????????????????? ???????????? ??????
		blockchain.Blockchain().Replace(payload)
	case MessageNewBlockNotify:
		// todo: validate the hash, timestamp, Txs, signatures......
		fmt.Printf("Received New Born Block from %s\n", p.key)
		var payload *blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Blockchain().AddPeerBlock(payload)
	case MessageNewTxNotify:
		fmt.Printf("Received New Transaction from %s\n", p.key)
		var payload *blockchain.Tx
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Mempool().AddPeerTx(payload)
	case MessageNewPeerNotify:
		fmt.Printf("Received New Peer from %s\n", p.key)
		var payload string
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		fmt.Printf("I will now /ws upgrade %s\n", payload)
		parts := strings.Split(payload, ":")
		AddPeer(parts[0], parts[1], parts[2], false)
	}
}
