package block

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kamogelosekhukhune777/go-blockchain/utils"
)

const (
	MiningDifficulty = 3
	MinningSender    = "THE BLOCKCHAIN"
	MinningReward    = 1.0
)

type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*Transaction
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transactions

	return b
}

func (b *Block) Print() {
	fmt.Printf("timestamp     %d\n", b.timestamp)
	fmt.Printf("nonce         %d\n", b.nonce)
	fmt.Printf("previoushash  %d\n", b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}

func (b *Block) MarshalJson() ([]byte, error) {
	return json.Marshal(struct {
		Nonce        int            `json:"nonce"`
		PreviousHash [32]byte       `json:"previous_hash"`
		Timestamp    int64          `json:"time_stamp"`
		Transaction  []*Transaction `json:"transaction"`
	}{
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Timestamp:    b.timestamp,
		Transaction:  b.transactions,
	})
}

type BlockChain struct {
	transactionpool   []*Transaction
	chain             []*Block
	blockChainaddress string
}

func NewBlockChain(blockChainaddress string) *BlockChain {
	b := &Block{}
	bc := new(BlockChain)
	bc.blockChainaddress = blockChainaddress
	bc.CreateBlock(0, b.Hash())

	return bc
}

func (bc *BlockChain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionpool)
	bc.chain = append(bc.chain, b)
	bc.transactionpool = []*Transaction{}

	return b
}

func (bc *BlockChain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *BlockChain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func (bc *BlockChain) AddTransaction(sender, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := NewTransactions(sender, recipient, value)
	if sender == MinningSender {
		bc.transactionpool = append(bc.transactionpool, t)
		return true
	}
	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		bc.transactionpool = append(bc.transactionpool, t)
		return true
	} else {
		log.Println("ERROR: verify Transaction")
	}
	return false
}
func (bc *BlockChain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

func (bc *BlockChain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionpool {
		transactions = append(transactions,
			NewTransactions(t.senderBlockChainAddress, t.recipientBlockChainAddress, t.value))
	}
	return transactions
}

func (bc *BlockChain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{0, nonce, previousHash, transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
}

func (bc *BlockChain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MiningDifficulty) {
		nonce += 1
	}
	return nonce
}

func (bc *BlockChain) Minning() bool {
	bc.AddTransaction(MinningSender, bc.blockChainaddress, MinningReward, nil, nil)
	nonce := bc.ProofOfWork()
	previousHah := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHah)
	log.Println("action=minning, status=success")
	return true
}

func (bc *BlockChain) CalculateTotalAmmount(blockChainAddress string) float32 {
	var totalAmount float32
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if blockChainAddress == t.recipientBlockChainAddress {
				totalAmount += value
			}
			if blockChainAddress == t.senderBlockChainAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

type Transaction struct {
	senderBlockChainAddress    string
	recipientBlockChainAddress string
	value                      float32
}

func NewTransactions(sender, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_blockchain_address      %s\n", t.senderBlockChainAddress)
	fmt.Printf(" recipient_blockchain_address   %s\n", t.recipientBlockChainAddress)
	fmt.Printf(" value                          %.1f\n", t.value)
}

func (t *Transaction) MarshalJson() ([]byte, error) {
	return json.Marshal(struct {
		SenderBlockChainAddress    string  `json:"sender_blockchain_Address"`
		RecipientBlockChainAddress string  `json:"recipient_Blockchain_Address"`
		Value                      float32 `json:"value"`
	}{
		SenderBlockChainAddress:    t.senderBlockChainAddress,
		RecipientBlockChainAddress: t.recipientBlockChainAddress,
		Value:                      t.value,
	})
}
