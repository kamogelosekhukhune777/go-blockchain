package main

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/kamogelosekhukhune777/go-blockchain/block"
	"github.com/kamogelosekhukhune777/go-blockchain/wallet"
)

var cache map[string]*block.BlockChain = make(map[string]*block.BlockChain)

type BlockchainServer struct {
	port uint16
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port: port}
}

func (bcs *BlockchainServer) Port() uint16 {
	return bcs.port
}

func (bcs *BlockchainServer) GetBlockChain() *block.BlockChain {
	bc, ok := cache["blockchain"]
	if !ok {
		minersWallet := wallet.NewWallet()
		bc = block.NewBlockChain(minersWallet.BlockChainAddress(), bcs.Port())
		cache["blockchain"] = bc
		log.Printf("private_key:%v", minersWallet.PrivateKeyStr())
		log.Printf("public_key:%v", minersWallet.PublicKeyStr())
		log.Printf("blockxhain_address:%v", minersWallet.BlockChainAddress())
	}
	return bc
}

func (bcs *BlockchainServer) GetChain(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-type", "application/json")
		bc := bcs.GetBlockChain()
		m, _ := bc.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		log.Printf("ERROR: Invalid HTTP Method ")
	}
}

func (bcs *BlockchainServer) Run() {
	http.HandleFunc("/", bcs.GetChain)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcs.Port())), nil))
}
