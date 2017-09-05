package blockchain

import (
	"fmt"
	"crypto/sha256"
	"math"
	"time"
)

type MerkleTree []byte

type Block struct{
	index int
	transactions []Transaction
	txs MerkleTree
	prev_hash []byte
	timestamp int64
	nonce []byte
	hash []byte
}

func NewBlock(i int, tr []Transaction, prhash prev_hash, n []byte]) *Block{
	bolake:=Block{index: i, transactions: tr, prev_hash: prhash, timestamp: time.Now().UnixNano(), nonce: n}
	bolake.GenerateMTree()
	return bolake
}
func (b *Block) GenesisBlock() *Block{
	return var Block{nonce: []byte{108, 53, 181, 0}}
}
func (b *Block) Hash() []byte {
	mix:=[]byte(fmt.Sprint(b.index, b.txs, b.prev_hash, b.timestamp, b.nonce))
	hash:=sha256.Sum256(mix)
	return hash[:]
}
func (b *Block) ValidateIndex(prev_block Block) bool {
	return prev_block.index==b.index-1
}
func (b *Block) ValidatePrevHash(prev_block Block) bool {
	return prev_block.hash==b.prev_hash
}
func (b *Block) ValidateHash() bool {
	return b.hash==b.Hash()
}
func (b *Block) ValidateNonce(difficulty int) bool {
	res:=make([]byte, difficulty)
	return sha256.Sum256(b.nonce)[0:difficulty]==res
}
func (b *Block) GenerateMTree(){
	ts:=[][]byte(b.transactions)
	op:=math.Log2(float64(len(ts)))
	dif:=op-float64(int(op))
	var N int
	if dif==0{
		N=int(op)
	}else{
		N=int(op)+1
	}
	for j:=0; j<N; j++{
		var L int
		if len(ts)%2==0{
			L=len(ts)/2
		}else{
			L=int(len(ts)/2)+1
		}
		layer:=make([][]byte, L, 32)
		for i:=0; i<len(ts); i+=2{
			var ix []byte
			ix=append(ix, ts[i]...)
			if i!=len(ts)-1{
				ix=append(ix, ts[i+1]...)
			}
			hash:=sha256.Sum256(ix)
			layer[i/2]=hash[:]
		}
		//fmt.Println(">>",layer)
		ts=layer
		if len(ts)==1{
			b.txs=ts[0]
		}
	}
}