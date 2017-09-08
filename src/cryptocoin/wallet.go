package cryptocoin

import ("crypto/rand"
		"encoding/json"
		"os"
		"crypto/rsa"
		"fmt"
		)
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
	var txs map[int]Transaction
	JaceUp(txtname, &txs)
	return txs
}
func GenerateTxTable(txtname string){
	txs:=map[int]Transaction{}
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