package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"strconv"
	"time"
)

type Block struct {
	// Block header
	Timestamp     int64
	PrevBlockHash []byte // hash of previous block
	Hash          []byte // hash of itself
	Nonce         int

	// Block body
	Transactions []*Transaction
}

// NewBlock creates and returns Block
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := new(Block)
	block.PrevBlockHash = prevBlockHash
	block.Timestamp = time.Now().Unix()
	block.Transactions = transactions

	pow := NewProofOfWork(block, TargetBits)
	block.Nonce, block.Hash = pow.Run()
	return block
}

// NewGenesisBlock creates and returns genesis block
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

func (b *Block) Print() {
	fmt.Printf("Hash: %x\n", b.Hash)
	fmt.Printf("PrevBlockHash: %x\n", b.PrevBlockHash)
	fmt.Println("Timestamp: ", b.Timestamp)
	fmt.Printf("PoW: %s\n", strconv.FormatBool(NewProofOfWork(b, TargetBits).Validate()))
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(b); err != nil {
		log.Panic(err)
	}
	return result.Bytes()

}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	hash := sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return hash[:]
}

func (b *Block) IsGenesisBlock() bool {
	return len(b.PrevBlockHash) == 0
}

func DeserializeBlock(d []byte) *Block {
	block := new(Block)

	decoder := gob.NewDecoder(bytes.NewReader(d))
	if err := decoder.Decode(block); err != nil {
		log.Panic(err)
	}

	return block

}
