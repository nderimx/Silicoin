package cryptocoin

import (
	"fmt"
	"crypto/sha256"
	"crypto/rsa"
	"crypto/rand"
	"time"
	"crypto"
)

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
	for i:=0; i<len(txs); i++{
		ttxs:=[]Transaction{txs[i]}
		for j:=i+1; j<len(txs); j++{
			if bytes.Equal(txs[i].prev_location.trans_hash, txs[j].prev_location.trans_hash){
				ttxs=append(ttxs, txs[j])
				txs=append(txs[:j], txs[j+1:]...)
			}
		}
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
			prblcki:=ttxs[0].prev_location.block_index
			cblck1:=bc.LatestBlock().index //ad-hoc for now
			for j:=prblcki+1; j<=cblck1; j++{
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
			btx:=bc.blocks[prblcki].transactions
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
