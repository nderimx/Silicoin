package node

import (
	"fmt"
	"math"
)

const Port int=8882
const IntroNodes []inetAdress={filan,fistek,hasan,hysen}
var ID int;

/*
type IPBoard struct {
	//inetAddress still pseudo
	address inetAddress
	port int
}*/

type Node struct {
	//later make it an actual hash table
	node_table map[int]inetAddress
	wallet Wallet
	collected_transactions []Transaction
	block_chain BlockChain
}
func (n *Node) gather_transactions() []Transaction {
	//all javaesque pseudo code from here on out
	transactions:=[9]Transaction
	packet:=make(DataGramPacket{empty})
	sock:=make(Socket{8889})
	sock.receive(packet)
	index:=0
	for sameJsonFormat(Transaction{}.to_json, packet)||index<=9{
		new_transaction Transaction=json_unpack(packet)
		if new_transaction.verify(){
			transactions[index++]=new_transaction
			n.broadcast_transaction(new_transaction)
		}
		packet:=make(DataGramPacket{empty})
		sock:=make(Socket{8889})
		sock.receive(packet)
	}
	return n.collected_transactions=transactions
}
func (n *Node) broadcast_transaction(transaction Transaction) bool {
	for peer:=range n.node_table{
		packet:=make(DataGramPacket{transaction.to_json})
		sock:=make(Socket{peer.port})
		sock.send(peer.address, packet)
	}
}
func (n *Node) broadcast_block(block Block) bool {
	for peer:=range n.node_table{
		packet:=make(DataGramPacket{block.to_json})
		sock:=make(Socket{peer.port})
		sock.send(peer.address, packet)
	}
}
func (n *Node) receive_new_block(bc BlockChain){

	packet:=make(DataGramPacket{empty})
	sock:=make(Socket{8889})
	sock.receive(packet)
	index:=0
	for sameJsonFormat(Block{}.to_json, packet){
		new_block Block=json_unpack(packet)
		if bc.receive_block(new_block){
			n.broadcast_block(new_block)
		}
		packet:=make(DataGramPacket{empty})
		sock:=make(Socket{packet})
		sock.receive(8889)
	}
}
func (n *Node) request_block_chain() bool {
//later
}
func (n *Node) send_block_chain(block_chain BlockChain) bool {
//later
}
func (n *Node) get_random_intro() inetAddress{
	i:=math.random.Seed(0,3)
	return IntroNodes[i]
}
func (n *Node) join_net(intro inetAddress){
	join int=8421
	packet:=make(DataGramPacket{})
	sock:=make(Socket{join})
	sock.send(intro, packet)

	packet:=make(DataGramPacket{empty})
	sock:=make(Socket{join})
	sock.receive(packet)
	ID=packet.GetData()
}