package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/base58"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockChainAddress string
}

func NewWallet() *Wallet {
	w := new(Wallet)
	//creating ECDSA private key (32 bytes) public key(64 bytes)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.privateKey = privateKey
	w.publicKey = &w.privateKey.PublicKey
	//perform SHA-256 hashing on public key(32 bytes).
	h2 := sha256.New()
	h2.Write(w.publicKey.X.Bytes())
	h2.Write(w.publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)
	//perform RIPEMD-160 hashing on the result of SHA-256(20 bytes)

	//sha-256
	h3 := sha256.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)
	//add version byte in front of RIPEMD-160 hash(0x00 for main Network)
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])
	//perform SHA-256 hash on the extended RIPEMD-160 result
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)
	//perform SHA-256 hash on the result of the previous SHA-256 hash
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)
	//take the first 4 bytes of the second SHA-256 hash for checksum
	checksum := digest6[:4]
	//add the 4 checksum bytes from 7 at the end 0f extended RIPEMD-160 hash from 4 (25 bytes)
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[:21], checksum[:])
	//convert the result from a byte string into base58
	address := base58.Encode(dc8)
	w.blockChainAddress = address

	return w
}

func (w *Wallet) PrivateKey() ecdsa.PrivateKey {
	return *w.privateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() ecdsa.PublicKey {
	return *w.publicKey
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) BlockChainAddress() string {
	return w.blockChainAddress
}

//https://github.com/btcsuite/btcd/btcutil/base58/
