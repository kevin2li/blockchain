package pkg

import (
	"fmt"
	"log"
)

func MultiInputHeuristic(addr string, tx Transaction) []string {
	if IsInTxInputs(addr, tx) {
		in_addrs := GetTxInAddrs(tx)
		return in_addrs
	}
	return nil
}

func CoinbaseHeuristic(addr string, tx Transaction) []string {
	if IsCoinbaseTx(tx) && IsInTxOutputs(addr, tx) {
		out_addrs := GetTxOutAddrs(tx)
		return out_addrs
	}
	return nil
}

// TODO: ChangeHeuristic
func ChangeHeuristic(addr string, tx Transaction) []string {

	return nil
}

func ClusterByAddr(addr string, txs []Transaction, addrList chan []string) {
	var result []string
	result = append(result, addr)
	for _, tx := range txs {
		// rule1
		out := MultiInputHeuristic(addr, tx)
		if out != nil {
			result = append(result, out...)
		}
		// rule2
		out = CoinbaseHeuristic(addr, tx)
		if out != nil {
			result = append(result, out...)
		}
		// rule3
		out = ChangeHeuristic(addr, tx)
		if out != nil {
			result = append(result, out...)
		}
	}
	result = Unique(result)
	addrList <- result
}

func Cluster(addr string, txs []Transaction) []string {
	var finalAddrList = make(HashSet)
	finalAddrList.Add(addr)
	var queue = make([]string, 0)
	queue = append(queue, addr)
	var iterations = 1
iter:
	fmt.Printf("================================Iteration %d started!================================\n", iterations)
	fmt.Printf("INFO: total: %d addresses.\n", len(queue))
	var n = len(queue)
	addrList := make(chan []string, n)
	for i := 0; i < n; i++ {
		addr := queue[i]
		fmt.Printf("[%d/%d] Starting cluster from address: %s\n", i+1, n, addr)
		go ClusterByAddr(addr, txs, addrList)
	}
	queue = make([]string, 0) // clear queue
	for i := 0; i < n; i++ {
		addrs := <-addrList
		for _, addr := range addrs {
			if _, ok := finalAddrList[addr]; !ok {
				queue = append(queue, addr)
				finalAddrList.Add(addr)
			}
		}
	}
	// whether have new address
	if len(queue) > 0 {
		fmt.Printf("INFO: new addresses added: %+v\n\n", queue)
		iterations++
		goto iter
	}
	result := finalAddrList.GetData()
	return result
}

func StartCluster(dataset_path string, start_addr string) {
	fmt.Println("INFO: Loading transactions....")
	all_txs, err := ReadTransaction(dataset_path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("INFO: Load total %d transactions done!\n", len(all_txs))
	fmt.Printf("INFO: Start cluster from address: %s!\n", start_addr)
	result := Cluster(start_addr, all_txs)
	fmt.Println("\n--------------------------------Cluster Finished!--------------------------------")
	fmt.Printf("INFO: cluster total %d addresses, final cluster result:\n %v\n", len(result), result)
}
