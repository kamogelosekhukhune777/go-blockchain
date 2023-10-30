package main

import (
	"fmt"
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
func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	fmt.Printf("test")
}
