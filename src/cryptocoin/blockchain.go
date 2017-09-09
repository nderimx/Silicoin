package cryptocoin

import (
	"crypto/sha256"
	"bytes"
)

type BlockChain struct{
	blocks []Block
}

func NewChain(prk rsa.PrivateKey, pbk rsa.PublicKey) BlockChain{
	nonce:=GenerateNonce([]byte("GENESIS"), 3)
	now:=time.Now().UnixNano()
	h:=sha256.Sum256([]byte(fmt.Sprint(pbk)))
	reward:=Transaction{
						pub_key: pbk,
	 					amount: 50,
	 					timestamp: now,
						hash: h[:]}
	reward.Sign(prk)
	return BlockChain{[]Block{NewBlock(0, []Transaction{Transaction{}}, reward, []byte("GENESIS"), nonce)}}
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
func (bc *BlockChain) GenerateBlock(transactions []Transaction, reward Transaction) Block{
	prev:=bc.LatestBlock()
	prvhash:=prev.Hash()
	bolake:=NewBlock(prev.index+1, transactions, reward, prvhash, GenerateNonce(prvhash, bc.GetDifficulty()))
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
		//validate all txs
		return nil
	}
	return errors.New("Block isn't valid!")
}
func (bc *BlockChain) GetDifficulty() int{
	//later relate the number to the avg time it takes to generate a block
	return 3
}
func (bc *BlockChain) GetRewardAmount() float64{
	//bc.LatestBlock().index  << relate to chain length
	return 50
}
func (bc *BlockChain) Equals(fbc BlockChain) bool{
	return false
	//edit
}
func IsGoodChain(fbc BlockChain) bool{
	return false
	//edit
}