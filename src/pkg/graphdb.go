package pkg

import (
	"fmt"
	"log"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/pkg/errors"
)

type Params map[string]interface{}

func Filter(txs []Transaction, addrs []string) []Transaction {
	result := make([]Transaction, 0)
	set := HashSet{}
	set.Init(addrs)
	for _, tx := range txs {
		all_addrs := append(GetTxInAddrs(tx), GetTxOutAddrs(tx)...)
		for _, a := range all_addrs {
			if set[a] {
				result = append(result, tx)
				break
			}
		}
	}
	return result
}

func InsertTransaction(driver neo4j.Driver, tx Transaction) error {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	in_addrs, out_addrs := GetTxInAddrs(tx), GetTxOutAddrs(tx)
	var createTx_cql = "MERGE (tx:Transaction {txid: $txid, in_degree: $in_degree, out_degree: $out_degree, time: $time, height: $height})"
	// 1. create tx node
	params := Params{
		"txid":       tx.Txid,
		"in_degree":  len(in_addrs),
		"out_degree": len(out_addrs),
		"time":       GetTxTime(tx),
		"height":     tx.Block.Height,
	}
	var insertInputFn = func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run(createTx_cql, params)
		if err != nil {
			return nil, err
		}
		return records, nil
	}
	_, err := session.WriteTransaction(insertInputFn)
	if err != nil {
		err = errors.Wrap(err, "insert transaction failed!")
		return err
	}
	// 2. create tx input addrs node
	for _, addr := range in_addrs {
		params := Params{
			"address1": addr,
			"txid":     tx.Txid,
		}
		// create addr node
		var insertInputFn = func(tx neo4j.Transaction) (interface{}, error) {
			records, err := tx.Run("MERGE (addr1:Addr { address: $address1 }) RETURN addr1", params)
			if err != nil {
				return nil, err
			}
			return records, nil
		}
		_, err := session.WriteTransaction(insertInputFn)
		if err != nil {
			err = errors.Wrap(err, "insert transaction failed!")
			return err
		}
		// create relationship
		insertInputFn = func(tx neo4j.Transaction) (interface{}, error) {
			records, err := tx.Run("MATCH (addr1:Addr { address: $address1 }), (tx:Transaction {txid: $txid}) CREATE (addr1)-[:In]->(tx)  RETURN addr1, tx", params)
			if err != nil {
				return nil, err
			}
			return records, nil
		}
		_, err = session.WriteTransaction(insertInputFn)
		if err != nil {
			err = errors.Wrap(err, "insert transaction failed!")
			return err
		}
	}
	// 3. create tx output addrs node
	for _, addr := range out_addrs {
		params := Params{
			"address2": addr,
			"txid":     tx.Txid,
		}
		// create addr node
		var insertOutputFn = func(tx neo4j.Transaction) (interface{}, error) {
			records, err := tx.Run("MERGE (addr2:Addr { address: $address2 }) RETURN addr2", params)
			if err != nil {
				return nil, err
			}
			return records, nil
		}
		_, err := session.WriteTransaction(insertOutputFn)
		if err != nil {
			err = errors.Wrap(err, "insert transaction failed!")
			return err
		}
		// create relationship
		insertOutputFn = func(tx neo4j.Transaction) (interface{}, error) {
			records, err := tx.Run("MATCH (addr2:Addr { address: $address2 }), (tx:Transaction {txid: $txid }) CREATE (tx)-[:Out]->(addr2) RETURN addr2, tx", params)
			if err != nil {
				return nil, err
			}
			return records, nil
		}
		_, err = session.WriteTransaction(insertOutputFn)
		if err != nil {
			err = errors.Wrap(err, "insert transaction failed!")
			return err
		}
	}
	return nil
}

func Visualize(all_txs []Transaction, addrs []string) error {
	dbUri := "neo4j://localhost:7687"
	driver, err := neo4j.NewDriver(dbUri, neo4j.BasicAuth("neo4j", "test", ""))
	if err != nil {
		err = errors.Wrap(err, "")
		return err
	}
	defer driver.Close()
	txs := Filter(all_txs, addrs)
	var n = len(txs)
	bar := GetProgressBar(n)
	defer bar.Close()
	for _, tx := range txs {
		err = InsertTransaction(driver, tx)
		if err != nil {
			err = errors.Wrap(err, "")
			return err
		}
		bar.Add(1)
	}
	fmt.Println("Done!")
	return nil
}

func StartViz(dataset_path string, addr_path string) error {
	fmt.Println("INFO: Loading transactions....")
	all_txs, err := ReadTransaction(dataset_path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("INFO: Load total %d transactions done!\n", len(all_txs))
	fmt.Println("INFO: Reading addresses....")
	addrs, err := ReadAddrs(addr_path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("INFO: Read total %d addresses done!\n", len(addrs))
	fmt.Println("INFO: Start insert transactions into neo4j...")
	err = Visualize(all_txs, addrs)
	if err != nil {
		err = errors.Wrap(err, "")
		return err
	}
	fmt.Println("INFO: Insert transactions done...")
	return nil
}
