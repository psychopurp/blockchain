package main

import "fmt"

type BlockChain struct {
	blocks []*Block
}

func NewBlockChain() *BlockChain {
	bc := new(BlockChain)
	bc.blocks = []*Block{NewGenesisBlock()}
	return bc

}

func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

func (bc *BlockChain) Print() {
	for _, block := range bc.blocks {
		block.Print()
		fmt.Println()
	}
}
