package main

import "math"

var (
	TargetBits int = 16 // difficulty of Hash mining
	MaxNonce       = math.MaxInt64
	Subsidy    int = 10
)

const (
	dbFile              = "blockchain.db"
	blocksBucket        = "blocks"
	lastBlockKey        = "l"
	genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)
