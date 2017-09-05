package blockchain

import (
	"fmt"
	"math"
	"container/list"
	"time"
	"crypto/sha256"
	"bytes"
)

type BlockChain []Block

func NewChain() *BlockChain{
	return BlockChain{GenesisBlock()}
}
func (bc *BlockChain) GenerateBlock(transactions []Transaction) Block{
	prev:=bc.LatestBlock()
	bolake:=NewBlock(prev.index+1, transactions, prev.hash, bc.GenerateNonce(bc.GetDifficulty()))
	bolake.hash=bolake.Hash()
	append(bc, bolake)
	return bolake
}
func (bc *BlockChain) LatestBlock() Block{
	return bc.blocks[len(bc.blocks)]
}
func (bc *BlockChain) ReceiveBlock(block Block) bool {
	last_block:=bc.LatestBlock()
	if (
	block.ValidateIndex(last_block)&&
	block.ValidatePrevHash(last_block)&&
	block.ValidateHash()&&
	block.ValidateNonce(bc.GetDifficulty)
	){
		bc.blocks.add(block)
		return true
	}
	return false
}
func (bc *BlockChain) GetDifficulty() int {
	//later relate the number to the time it takes to generate a block
	return 3
}
func GenerateNonce(difficulty int) []byte {
	res:=make([]byte, difficulty)
	const basem=byte(255)
	nonce:=make([]byte, 1)
	for j:=0; ; j++{

		for i:=0; i<len(nonce); i++{
			if i!=0 && nonce[i]!=basem{
				nonce[i]++
				ph:=sha256.Sum256(nonce)
				if bytes.Equal(ph[0:difficulty], res){
					return nonce
				}
				i=-1
			}else if i!=0 && nonce[i]==basem{
				nonce[i]=byte(0)
			}else{
				for nonce[i]!=basem{
					nonce[i]++
					ph:=sha256.Sum256(nonce)
					if bytes.Equal(ph[0:difficulty], res){
						return nonce
					}
				}
				nonce[i]=byte(0)
			}
		}

		if j==len(nonce)-1{
			for k:=0; k<len(nonce); k++{
				nonce[k]=byte(0)
			}
			nonce=append(nonce, byte(0))
		}
	}
}