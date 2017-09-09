package cryptocoin

import (
	"crypto/sha256"
	"time"
	"bytes"
)
//later replace []Transactions with a struct that cointains MerkleRoot and it's full array of transactions
type Block struct{
	index int
	transactions []Transaction
	reward Transaction
	prev_hash []byte
	timestamp int64
	nonce []byte
}

func NewBlock(i int, tr []Transaction, rt Transaction, prhash []byte, n []byte) Block{
	bolake:=Block{index: i, transactions: tr, reward: rt, prev_hash: prhash, timestamp: time.Now().UnixNano(), nonce: n}
	return bolake
}
func (b *Block) Hash() []byte {
	mix:=[]byte(fmt.Sprint(b.transactions, b.reward, b.prev_hash, b.timestamp, b.nonce))
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
	h:=sha256.Sum256(append(b.nonce, b.prev_hash...))
	return bytes.Equal(h[0:difficulty], res)
}
func GenerateNonce(prev_hash []byte, difficulty int) []byte{
	res:=make([]byte, difficulty)
	const basem=byte(255)
	nonce:=make([]byte, 1)
	for j:=0; ; j++{

		for i:=0; i<len(nonce); i++{
			if i!=0 && nonce[i]!=basem{
				nonce[i]++
				ph:=sha256.Sum256(append(nonce, prev_hash...))
				if bytes.Equal(ph[0:difficulty], res){
					return nonce
				}
				i=-1
			}else if i!=0 && nonce[i]==basem{
				nonce[i]=byte(0)
			}else{
				for nonce[i]!=basem{
					nonce[i]++
					ph:=sha256.Sum256(append(nonce, prev_hash...))
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