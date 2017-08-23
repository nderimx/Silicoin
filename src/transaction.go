package transaction

import (
	"fmt"
	"math"
)

type TransactionLocation struct{
	block_index int
	index int
}
type Transaction{
	index int
	pub_key []byte
	amount float64
	trans_hash []byte
	prev_location TransactionLocation int //location of the block the previous transaction is embedded in... and the transaction
	signature []byte
	timestamp int
}
//JsonObject is a pseudo object
func (t *Transaction) to_json() JsonObject {
	//just find a legin json library
}
func (t *Transaction) gen_hash(prev_hash []byte){
	//not sure how to combine just yet
	mix []byte=t.pub_key+prev_hash+TOBYTES(t.amount)+TOBYTES(t.timestamp)+TOBYTES(t.prev_location)
	//find a sha256 library
	t.trans_hash=hash(mix)
}
func (t *Transaction) sign(prk []byte){
	//find asymetric key encryption library
	t.signature=encrypt(t.trans_hash, prk)
}
func (t *Transaction) verify(prev_pub_key []byte) bool{
	//find asymetric key encryption library
	return decrypt(t.signature, prev_pub_key)==t.trans_hash
}
func (t *Transaction) find_prev_key(bc BlockChain) []byte{
	//assuming that's how i'm gonna call a location inside a blockchain
	block:=bc[t.prev_location.block_index]
	pt:=block.transactions[t.prev_location.index]
	return pt.pub_key
}