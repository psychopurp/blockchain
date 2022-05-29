package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlock_Serialize(t *testing.T) {
	b := &Block{}

	d := b.Serialize()

	assert.Equal(t, b, DeserializeBlock(d))
}
