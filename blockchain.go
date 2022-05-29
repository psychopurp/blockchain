package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	bolt "go.etcd.io/bbolt"
)

type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

// NewBlockchain creates a new Blockchain with genesis Block
func NewBlockChain(address string) *BlockChain {
	if !dbExist() {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte(lastBlockKey))
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := &BlockChain{tip: tip, db: db}
	return bc

}

func CreateBlockChain(address string) *BlockChain {
	if dbExist() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)

		err = bucket.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = bucket.Put([]byte(lastBlockKey), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}

		tip = genesis.Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &BlockChain{db: db, tip: tip}
}

// MineBlock mines a new block with the provided transactions
func (bc *BlockChain) MineBlock(transactions []*Transaction) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte(lastBlockKey))
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = newBlock.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

// FindUTXO find all unspent transaction outputs
func (bc *BlockChain) FindUTXO() map[string]*UTXOutputs {

	references := make(map[string][]int)
	UTXOs := make(map[string]*UTXOutputs)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		findUTXO(block.Transactions, references, UTXOs)

		if block.IsGenesisBlock() {
			break
		}

	}

	return UTXOs
}

func (bc *BlockChain) FindSpendableUTXOs(address string, amount int) (int, map[string][]int) {
	UTXOs := bc.FindUTXO()
	unspentOutputs := make(map[string][]int)
	accumulated := 0

WORK:
	for txID, outs := range UTXOs {
		for idx, out := range outs.Outputs {
			if out.CanUnlockOutputWith(address) {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outs.Indexes[idx])

				if accumulated >= amount {
					break WORK
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{db: bc.db, currentHash: bc.tip}
}

func dbExist() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func findUTXO(transactions []*Transaction, references map[string][]int, UTXOs map[string]*UTXOutputs) {

	for _, tx := range transactions {
		outputs := new(UTXOutputs)
		txID := hex.EncodeToString(tx.ID)

	OUTPUTS:
		for outIdx, out := range tx.Vout {

			for _, usedIdx := range references[txID] {
				if outIdx == usedIdx {
					continue OUTPUTS
				}
			}

			outputs.Outputs = append(outputs.Outputs, out)
			outputs.Indexes = append(outputs.Indexes, outIdx)
		}

		UTXOs[txID] = outputs

		if !tx.IsCoinbase() {
			for _, in := range tx.Vin {
				inTXID := hex.EncodeToString(in.TxID)
				referenceID := in.Vout
				references[inTXID] = append(references[inTXID], referenceID)
			}
		}

	}

}
