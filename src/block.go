package block

import (
	"fmt"
	"math"
)

type Block struct{
	index int
	transactions []Transaction
	prev_hash []byte
	timestamp int
	nonce []byte
	hash []byte
}
func (b *Block) genesis_block() Block{
	//return b:=Block{}
}
func (b *Block) to_hash() []byte {
	//not sure how to combine just yet
	mix []byte=b.index+b.transactions+b.prev_hash+b.timestamp+b.nonce
	//find a sha256 library
	hash(mix)
}
//JsonObject is a pseudo object
func (b *Block) to_json() JsonObject {
	//just find a legin json library
}
func (b *Block) validate_index(prev_block Block) bool {
	return prev_block.index==b.index-1
}
func (b *Block) validate_prev_hash(prev_block Block) bool {
	return prev_block.hash==b.prev_hash
}
func (b *Block) validate_hash() bool {
	return b.hash==b.to_hash()
}
func (b *Block) validate_nonce(difficulty int) bool {
	res [difficulty]byte
	for i:0;i<difficulty;i++{
		res[i]=0
	}
	return hash(b.nonce)[0:difficulty]==res
}