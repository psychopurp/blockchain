package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
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
	Data []byte
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := new(Block)
	block.PrevBlockHash = prevBlockHash
	block.Timestamp = time.Now().Unix()
	block.Data = []byte(data)

	pow := NewProofOfWork(block, TargetBits)
	block.Nonce, block.Hash = pow.Run()
	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func setHash(timestamp int64, prevBlockHash []byte, data []byte) []byte {
	headers := bytes.Join([][]byte{
		[]byte(strconv.FormatInt(timestamp, 10)),
		prevBlockHash,
		data,
	}, []byte{})

	hash := sha256.Sum256(headers)
	return hash[:]
}

func (b *Block) Print() {
	fmt.Printf("block : %s\n", b.Data)
	fmt.Printf("Hash: %x\n", b.Hash)
	fmt.Printf("PrevBlockHash: %x\n", b.PrevBlockHash)
	fmt.Println("Timestamp: ", b.Timestamp)
	fmt.Printf("PoW: %s\n", strconv.FormatBool(NewProofOfWork(b, TargetBits).Validate()))
}
