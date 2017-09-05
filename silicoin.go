package main

import (
	"fmt"
	"math"
	"io/ioutil"
	"encoding/json"
	"coin/blockchain"
	"coin/transaction"
	"coin/node"
)

func Init() (Node, error){
	nt, err := ioutil.ReadFile("node_table")
	if err==nil{
		var nodes map[int]string
		er:=json.Unmarshal(nt, &nodes)
		if er==nil{
			node:=Node{ID: 1, node_table: nodes, transaction_index:0}
			node.collected_transactions=make([]Transaction, 0, TRXARRAYSIZE)
			f, err:=ioutil.ReadFile("chain")
			check(err)
			node.block_chain=BlockChain(f)
			return node, nil
		}else {return Node{}, er}
	}else {return Node{}, err}
}
func main(){
	node, err:=Init()
	fmt.Println(err)
	//node.receive_blocks(node.block_chain)
	go node.gather_transactions()
	node.ui()
}












































