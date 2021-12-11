package main

import "blockchain/src/cmd"

func main() {
	cmd.Execute()
	// txs, err := pkg.ReadTransaction("/home/likai/code/go_program/go_learn2/result/block_height=711901.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// tx := txs[11]
	// in_addrs := pkg.GetTxInAddrs(tx)
	// out_addrs := pkg.GetTxOutAddrs(tx)
	// fmt.Printf("txid: %+v\n", tx.Txid)
	// fmt.Printf("in: %+v\n", in_addrs)
	// fmt.Printf("out: %+v\n", out_addrs)
}
