package cryptocoin

import (
	"crypto/sha256"
	"bytes"
)

type BlockChain struct{
	blocks []Block
}

func NewChain(prk rsa.PrivateKey, pbk rsa.PublicKey) BlockChain{
	nonce:=GenerateNonce(3)
	now:=time.Now().UnixNano()
	h:=sha256.Sum256([]byte(fmt.Sprint(pbk)))
	transaction:=Transaction{
						pub_key: pbk,
	 					amount: 60,
	 					timestamp: now,
						hash: h[:]}
	transaction.Sign(prk)
	transactions:=[]Transaction{transaction}
	return BlockChain{[]Block{NewBlock(0, transactions, []byte("GENESIS"), nonce)}}
}
func (bc *BlockChain) MinedTransaction(prk rsa.PrivateKey, pbk rsa.PublicKey) Transaction{
	ammo:=bc.GetRewardAmount()
	now:=time.Now().UnixNano()
	h:=sha256.Sum256([]byte(fmt.Sprint(pbk)))
	transaction:=Transaction{
						pub_key: pbk,
	 					amount: ammo,
	 					timestamp: now,
						hash: h[:]}
	transaction.Sign(prk)
	return transaction
}
func (bc *BlockChain) GenerateBlock(transactions []Transaction) Block{
	prev:=bc.LatestBlock()
	bolake:=NewBlock(prev.index+1, transactions, prev.Hash(), GenerateNonce(bc.GetDifficulty()))
	bc.blocks=append(bc.blocks, bolake)
	return bolake
}
func (bc *BlockChain) LatestBlock() Block{
	return bc.blocks[len(bc.blocks)-1]
}
func (bc *BlockChain) ReceiveBlock(block Block) error{
	last_block:=bc.LatestBlock()
	if (
	block.ValidateIndex(last_block)&&
	block.ValidatePrevHash(last_block)&&
	block.ValidateNonce(bc.GetDifficulty())){
		bc.blocks=append(bc.blocks, block)
		return nil
	}
	return errors.New("Block isn't valid!")
}
func (bc *BlockChain) GetDifficulty() int{
	//later relate the number to the time it takes to generate a block
	return 3
}
func (bc *BlockChain) GetRewardAmount() float64{
	//bc.LatestBlock().index  << relate to chain length
	return 50
}
func GenerateNonce(difficulty int) []byte{
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
func (bc *BlockChain) Equals(fbc BlockChain) bool{
	return false
	//jk
}
func IsGoodChain(fbc BlockChain) bool{
	return false
}