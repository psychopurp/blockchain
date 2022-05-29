package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

// Transaction represents a Bitcoin transaction
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// TXInput represents a transaction input
type TXInput struct {
	TxID      []byte
	Vout      int
	ScriptSig string
}

// TXOutput represents a transaction output
type TXOutput struct {
	Value        int
	ScriptPubKey string
}

type UTXOutputs struct {
	Outputs []TXOutput
	Indexes []int
}

// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to string, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := []TXInput{{TxID: []byte{}, Vout: -1, ScriptSig: data}}
	txout := []TXOutput{{Value: Subsidy, ScriptPubKey: to}}

	tx := &Transaction{ID: nil, Vin: txin, Vout: txout}
	tx.SetID()

	return tx
}

// NewUTXOTransaction creates a new transaction
func NewUTXOTransaction(from, to string, amount int, bc *BlockChain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := bc.FindSpendableUTXOs(from, amount)
	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			inputs = append(inputs, TXInput{TxID: txID, Vout: out, ScriptSig: from})
		}
	}

	// Build a list of outputs
	outputs = append(outputs, TXOutput{Value: amount, ScriptPubKey: to})
	if acc > amount {
		outputs = append(outputs, TXOutput{Value: acc - amount, ScriptPubKey: from}) // a change
	}

	tx := &Transaction{
		Vin:  inputs,
		Vout: outputs,
	}
	tx.SetID()

	return tx
}

// IsCoinbase checks whether the transaction is coinbase
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].TxID) == 0 && tx.Vin[0].Vout == -1

}

// CanUnlockOutputWith checks whether the address initiated the transaction
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

// CanUnlockOutputWith checks if the output can be unlocked with provided data
func (out *TXOutput) CanUnlockOutputWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

// SetID sets ID of a transaction
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}
