package main

import (
	"log"
	"time"
)

type Block struct {
	nonce        int
	previousHash string
	timestamp    int64
	transaction  []string
}

func NewBlock(nonce int, previousHash string) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash

	return b
}

type Blockchain struct {
	transactionpool []string
	chain           []*Block
}

func NewBlockchain() *Blockchain {
	bc := new(Blockchain)
	bc.CreateBlock(0, "Init Hash")

	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash string) *Block {
	b := NewBlock(nonce, previousHash)
	bc.chain = append(bc.chain, b)

	return b
}

func init() {
	log.SetPrefix("Blockchain: ")
}
