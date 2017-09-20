package main

import ("fmt"
	"time"
	"crypto/sha256"
	"bytes"
	"crypto/rsa"
	"crypto/rand"
	"crypto"
	"io/ioutil"
	"encoding/json"
	"errors"
	"strconv"
	)
func main(){

	wallet1:=InitWallet("pr.key", "pb.key", "tz")
	wallet2:=InitWallet("pr2.key", "pb2.key", "tz2")
	fmt.Println("about to gen new chain...")
	chain1:=NewChain(wallet1.private_key, wallet1.public_key)
	wallet1.SaveTx(chain1.LatestBlock())
	chain2:=chain1
	fmt.Println("copied chain1 to chain2")
	//w1 paying w2 25 units
	fmt.Println("making a transaction...")
	firstTX, firstCHNG, terr:=NewTransaction(wallet1.private_key,
						wallet1.public_key,
						wallet2.public_key,
						26, wallet1.transactions[0], 0)
	check(terr)
	mtrx:=chain2.MinedTransaction(wallet2.private_key, wallet2.public_key)
	txs:=[]Transaction{firstTX, firstCHNG}
	fmt.Println("generating new block on chain2...")
	block:=chain2.GenerateBlock(txs, mtrx)
	fmt.Println("receiving new Block on chain1")
	err:=chain1.ReceiveBlock(block)
	check(err)
	wallet1.SaveTx(chain1.LatestBlock())
	wallet2.SaveTx(chain2.LatestBlock())
	b1:=chain1.LatestBlock()
	b2:=chain2.LatestBlock()
	eq:=bytes.Equal(b1.Hash(), b2.Hash())
	fmt.Println(chain1, "\n",chain2)
	ft.Println("both chains have the same last block?: ", eq)

}
////////
//Below methods get imported from /go/src/cryptocoin, but this is just a proof of concept
//Disclaimer!: recursive transaction verification on block validation yet to be added
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
//////////
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
	block.ValidateNonce(bc.GetDifficulty())&&
	VerifyTXS(block.transactions, bc)==nil){
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
//////////
type TransactionLocation struct{
	block_index int
	trans_hash []byte
}
type Transaction struct{
	pub_key rsa.PublicKey
	amount float64
	timestamp int64
	hash []byte
	prev_location TransactionLocation //location of the block the previous transaction is embedded in... and the transaction
	signature []byte
}


func NewTransaction(prk rsa.PrivateKey, opbk rsa.PublicKey, rpbk rsa.PublicKey,
					 ammo float64, prev Transaction, blck_index int) (Transaction, Transaction, error){
	now:=time.Now().UnixNano()
	var err error=nil
	payment:=Transaction{
						pub_key: rpbk,
	 					amount: ammo,
	 					timestamp: now,
						hash: Hash(prev.TransactionHash(), rpbk),
						prev_location: TransactionLocation{blck_index, prev.TransactionHash()}}
	payment.Sign(prk)
	var change Transaction
	diff:=prev.amount-ammo
	if diff>0{
		change=Transaction{
							pub_key: opbk,
		 					amount: diff,
		 					timestamp: now,
							hash: Hash(prev.TransactionHash(), opbk),
							prev_location: TransactionLocation{blck_index, prev.TransactionHash()}}
		change.Sign(prk)
	}else{
		err=errors.New("You don't have enough units to make this payment")
		payment=Transaction{}
		change=Transaction{}
	}
	return payment, change, err

}
//This method should pull transactions from the wallet, depending on the amount to be transfered
// func (w *Wallet) NewTransaction(rpbk rsa.PublicKey, ammo float64) (*Transaction, *Transaction){
// 	now:=time.Now().UnixNano()
// 	payment:=Transaction{
// 						pub_key: rpbk,
// 	 					amount: ammo,
// 	 					timestamp: now,
// 						hash: Hash(prev.TransactionHash(), rpbk),
// 						prev_location: TransactionLocation{blck_index, prev.TransactionHash()}
// 					}
// 	payment.Sign(prk)
// 	diff:=prev.amount-ammo
// 	if diff>0{
// 		change:=Transaction{
// 							pub_key: opbk,
// 		 					amount: diff,
// 		 					timestamp: now,
// 							hash: Hash(prev.TransactionHash(), opbk),
// 							prev_location: TransactionLocation{blck_index, prev.TransactionHash()}
// 						}
// 		change.Sign(prk)
// 	}else{
// 		change:=Transaction{}
// 	}
// 	return payment, change

// }
func Hash(prev_thash []byte, pub_key rsa.PublicKey) []byte{
	mix:=[]byte(fmt.Sprint(prev_thash, pub_key))
	hash:=sha256.Sum256(mix)
	return hash[:]
}
func (t *Transaction) TransactionHash() []byte{
	mix:=[]byte(fmt.Sprint(t.pub_key, t.hash, t.signature, t.amount, t.timestamp))
	hash:=sha256.Sum256(mix)
	return hash[:]
}
func (t *Transaction) Sign(prk rsa.PrivateKey){
	reader:=rand.Reader
	var err error
	t.signature, err=rsa.SignPKCS1v15(reader, &prk, crypto.SHA256, t.hash)
	if err!=nil{
		fmt.Println("Signing Error: ", err)
	}
}
func (t *Transaction) VerifySig(prev_pub_key rsa.PublicKey) error{
	err:=rsa.VerifyPKCS1v15(&prev_pub_key, crypto.SHA256, t.hash, t.signature)
	if err!=nil{
		return errors.New("Transaction signature not Valid!")
	}
	return nil
}
func (t *Transaction) FindPrevTrx(bc *BlockChain) (Transaction, bool){
	block:=bc.blocks[t.prev_location.block_index]
	//ptrans:=block.transactions[t.prev_location.index]
	txs:=block.transactions
	rw:=block.reward
	if bytes.Equal(rw.TransactionHash(), t.prev_location.trans_hash){
		return rw, true
	}
	ln:=len(txs)
	for i:=0; i<ln; i++{
		var prev_trans Transaction=txs[i]
		if bytes.Equal(prev_trans.TransactionHash(), t.prev_location.trans_hash){
			return txs[i], false
		}
	}
	return Transaction{}, false
}
/*
func VerifyTX(t Transaction, bc BlockChain, tAmouont float64) error{
	ptx, isReward:=t.FindPrevTrx(bc)
	sum:=ptx.amount+t.amount
	if (t.VerifySig(ptx.pub_key)==nil){
		if isReward{

			if  sum==bc.GetRewardAmount(){
				return nil
			}
		}else{
			VerifyTX(ptx, bc, difference)
		}
	}
	return errors.New("Transaction not Valid!"), -1
}*/
func (t *Transaction) VerifyHash(ptx Transaction) error{
	if bytes.Equal(ptx.TransactionHash(), t.prev_location.trans_hash){
		return nil
	}
	return errors.New("Tx hashes don't match!")
}
func VerifyTXS(txs []Transaction, bc *BlockChain) error{
	var prevBlockI int
	var currentBlockI int
	for i:=0; i<len(txs); i++{
		ttxs:=[]Transaction{txs[i]}
		for j:=i+1; j<len(txs); j++{
			if bytes.Equal(txs[i].prev_location.trans_hash, txs[j].prev_location.trans_hash){
				ttxs=append(ttxs, txs[j])
				txs=append(txs[:j], txs[j+1:]...)
			}
		}
		isLatestBlock:=true
		for{
			var s float64=0
			ptx, isMTRX:=ttxs[0].FindPrevTrx(bc)
			for j:=0; j<len(ttxs); j++{
				currentTransaction:=ttxs[j]
				sigError:=currentTransaction.VerifySig(ptx.pub_key)
				hashError:=currentTransaction.VerifyHash(ptx)
				if sigError==nil && hashError==nil{
					s+=currentTransaction.amount
				}else{
					return errors.New("Transaction has no link to another previous transaction")
				}
			}
			if ptx.amount!=s{
				return errors.New("Transaction amounts do not match!")
			}
			if isLatestBlock{
				currentBlockI=bc.LatestBlock().index
				isLatestBlock=false
			}else {
				currentBlockI=prevBlockI
			}
			prevBlockI=ttxs[0].prev_location.block_index
			for j:=prevBlockI+1; j<=currentBlockI; j++{
				cbt:=bc.blocks[j].transactions
				for k:=0; k<len(cbt); k++{
					if bytes.Equal(cbt[k].prev_location.trans_hash, ptx.TransactionHash()){
						return errors.New("Transaction already spent.")
					}
				}
			}
			if isMTRX{
				isMTRX=false
				break
			}
			ttxs:=[]Transaction{ptx}
			btx:=bc.blocks[prevBlockI].transactions
			for j:=i+1; j<len(btx); j++{
				if bytes.Equal(ptx.prev_location.trans_hash, btx[j].prev_location.trans_hash){
					ttxs=append(ttxs, btx[j])
				}
			}
		}
	}
	return nil
}
func (t *Transaction) FindPrevPubKey(bc *BlockChain) rsa.PublicKey{
	ptx, _:=t.FindPrevTrx(bc)
	return ptx.pub_key
}
func (t *Transaction) FindPrevTHash(bc *BlockChain) []byte{
	prtx, _:=t.FindPrevTrx(bc)
	return prtx.TransactionHash()
}
//////////
type Wallet struct {
	public_key rsa.PublicKey
	private_key rsa.PrivateKey
	transactions map[int]Transaction //map[block_index]transaction
}
func InitWallet(prname string, pbname string, txtname string) Wallet{
	prk, pbk:=GetPair(prname, pbname)
	txs:=GetTxTable(txtname)
	return Wallet{pbk, prk, txs}
}
//compare your public key with incoming
//block's transactions' public keys and save matching transactions with it's merkle branches to wallet.
func (w *Wallet) SaveTx(b Block){
	wpk:=w.public_key
	btxs:=b.transactions
	if b.reward.pub_key==wpk{
		w.transactions[b.index]=b.reward
	}
	for i:=0; i<len(btxs); i++{
		if wpk==btxs[i].pub_key{
			//w.transactions[b.hash]=btxs[i] on next version also save the merkle branches with merkle root
			w.transactions[b.index]=btxs[i]
		}
	}
}
func CreateWallet(prname string, pbname string, txtname string){
	GeneratePair(2048, prname, pbname)
	GenerateTxTable(txtname)
}
func GetPair(prname string, pbname string) (rsa.PrivateKey, rsa.PublicKey){
	var prk rsa.PrivateKey
	JaceUp(prname, &prk)
	var pbk rsa.PublicKey
	JaceUp(pbname, &pbk)
	return prk, pbk
}
func GeneratePair(bitSize int, prname string, pbname string){
	reader:=rand.Reader
	key, err:=rsa.GenerateKey(reader, bitSize)
	check(err)
	publicKey:=key.PublicKey

	SaveJsn(prname, key)
	SaveJsn(pbname, publicKey)
}
func GetTxTable(txtname string) map[int]Transaction{
	var txs map[string]Transaction
	JaceUp(txtname, &txs)
	return ConvMap(txs)
}
func GenerateTxTable(txtname string){
	txs:=map[string]Transaction{}
	SaveJsn(txtname, txs)
}
func check(e error){
	if e!=nil{
		fmt.Print(e)
	}
}
func SaveJsn(fileName string, ob interface{}){
	jfile, err:=json.Marshal(ob)
	check(err)
	ioutil.WriteFile(fileName, jfile, 0644)
}
func JaceUp(name string, t interface{}){
	file, err:=ioutil.ReadFile(name)
	check(err)
	err=json.Unmarshal(file, t)
	check(err)
}
func ConvMap(input map[string]Transaction) map[int]Transaction{
	iLen:=len(input)
	integers:=make([]int, iLen)
	transacters:=make([]Transaction, iLen)
	output:=make(map[int]Transaction, iLen)
	i:=0
	for key, value:=range input{
		integers[i], _=strconv.Atoi(key)
		transacters[i]=value
		i++
	}
	for i:=0; i<iLen; i++{
		output[integers[i]]=transacters[i]
	}
	return output
}
