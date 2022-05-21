package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockChain_AddBlock(t *testing.T) {

	bc := NewBlockChain()

	bc.AddBlock("Send 1 BTC to cacts.eth")
	bc.AddBlock("Send 2 more BTC to cacts.eth")

	bc.Print()

	assert.Equal(t, 3, len(bc.blocks))

}
