package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type ProofOfWork struct {
	block      *Block
	target     *big.Int
	difficulty int // difficulty of mining
}

func NewProofOfWork(block *Block, difficulty int) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	pow := &ProofOfWork{block: block, target: target, difficulty: difficulty}
	return pow
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < MaxNonce {
		hash = prepareHashData(pow.block, pow.difficulty, nonce)
		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			return nonce, hash[:]
		}

		nonce++
	}

	return -1, nil
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	hash := prepareHashData(pow.block, pow.difficulty, pow.block.Nonce)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}

func prepareHashData(block *Block, difficulty int, nonce int) [32]byte {
	data := bytes.Join([][]byte{
		block.PrevBlockHash,
		block.Data,
		[]byte(fmt.Sprint(block.Timestamp)),
		[]byte(fmt.Sprint(difficulty)),
		[]byte(fmt.Sprint(nonce)),
	}, []byte{})

	return sha256.Sum256(data)
}
