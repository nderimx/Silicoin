package wallet
import ("crypto/rand"
		"encoding/gob"
		"os"
		"crypto/rsa"
		"fmt"
		"container/list"
		)
type Wallet struct {
	public_key rsa.PublicKey
	private_key rsa.PrivateKey
	transactions map[int]Transaction //map[block_index]transaction
}
type Pouch list.List

func InitWallet(prname string, pbname string, txtname string) *Wallet{
	pbk, prk:=GetPair(prname, pbname)
	txs:=GetTxTable(txtname)
	return Wallet{pbk, prk, txs}
}
//after node creates or receives a transaction with it's own public key
func (p *Pouch) AddTx(t Transaction){
	//p=append(p, t)
}
//after node receives and verifies, or generates a new block
func (w *Wallet) AddTransaction(b Block, p Pouch){
	w.transactions[b.index]=p[len(p)-1]
}
func CreateWallet(prname string, pbname string, txtname string){
	GeneratePair(2048, prname, pbname)
	GenerateTxTable(txtname)
}
func GetPair(prname string, pbname string) (rsa.PrivateKey, rsa.PublicKey){
	var prk rsa.PrivateKey
	GobbleUp(prname, &prk)
	var pbk rsa.PublicKey
	GobbleUp(pbname, &pbk)
	return prk, pbk
}
func GeneratePair(bitSize int, prname string, pbname string){
	reader:=rand.Reader
	key, err:=rsa.GenerateKey(reader, bitSize)
	check(err)
	publicKey:=key.PublicKey

	SaveKey(prname, key)
	SaveKey(pbname, publicKey)
}
func GetTxTable(txtname string) map[int]Transaction{
	var txs map[int]Transaction
	GobbleUp(txtname, &txs)
	return txs
}
func GenerateTxTable(txtname string){
	reader:=rand.Reader
	txs:=map[int]Transaction{}
	SaveKey(txtname, txs)
}
func check(e error){
	if e!=nil{
		fmt.Print(e)
	}
}
func SaveKey(fileName string, key interface{}){
		outFile, err:=os.Create(fileName)
		check(err)
		defer outFile.Close()
		encoder:=gob.NewEncoder(outFile)
		err=encoder.Encode(key)
		check(err)
}
func GobbleUp(name string, t interface{}){
	file, err:=os.Open(name)
	check(err)
	defer file.Close()
	dec:=gob.NewDecoder(file)
	err=dec.Decode(t)
	check(err)
}