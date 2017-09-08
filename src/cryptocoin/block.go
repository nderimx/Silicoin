package cryptocoin

import (
	"crypto/sha256"
	"time"
	"bytes"
)
//later replace []Transactions with full Merkle Tree with the corresponding full Transactions
type Block struct{
	index int
	transactions []Transaction
	prev_hash []byte
	timestamp int64
	nonce []byte
}

func NewBlock(i int, tr []Transaction, prhash []byte, n []byte) Block{
	bolake:=Block{index: i, transactions: tr, prev_hash: prhash, timestamp: time.Now().UnixNano(), nonce: n}
	return bolake
}
func (b *Block) Hash() []byte {
	mix:=[]byte(fmt.Sprint(b.transactions, b.prev_hash, b.timestamp, b.nonce))
	hash:=sha256.Sum256(mix)
	return hash[:]
}
func (b *Block) ValidateIndex(prev_block Block) bool {
	return prev_block.index==b.index-1
}
func (b *Block) ValidatePrevHash(prev_block Block) bool {
	return bytes.Equal(prev_block.Hash(), b.prev_hash)
}
func (b *Block) ValidateNonce(difficulty int) bool {
	res:=make([]byte, difficulty)
	h:=sha256.Sum256(b.nonce)
	return bytes.Equal(h[0:difficulty], res)
}