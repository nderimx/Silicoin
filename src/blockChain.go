package blockChain

import (
	"fmt"
	"math"
)

type BlockChain struct {
	//list<T> is a pseudo object
	blocks List<Block>
}
func (bc *BlockChain) latest_block() Block {
	return bc.blocks<bc.blocks.size()>
}
func (bc *BlockChain) generate_block(transactions []Transaction) Block {
	prev:=bc.latest_block()
	block:=Block{
		prev.index+1, transactions, prev.hash, getTime(), bc.generate_nonce(bc.get_difficulty())
	}
}
func (bc *BlockChain) receive_block(block Block) bool {
	last_block:=bc.latest_block()
	if (
	block.validate_index(last_block)&&
	block.validate_prev_hash(last_block)&&
	block.validate_hash()&&
	block.validate_nonce(bc.get_difficulty)
	){
		bc.blocks.add(block)
		return true
	}
	return false
}
func (bc *BlockChain) get_difficulty() int {
	//later relate the number to the time it takes to generate a block
	return 5
}
func (bc *BlockChain) generate_nonce(difficulty int) []byte {
	res [difficulty]byte
	for i:0;i<difficulty;i++{
		res[i]=0
	}
	//get the correct []byte notations/syntax
	nonce []byte=0
	for i:=0;hash(i)[0:difficulty]!=res;i++{
		nonce:=i
	}
	return nonce
}