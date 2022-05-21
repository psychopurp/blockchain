package main

func main() {
	bc := NewBlockChain()

	bc.AddBlock("Send 1 BTC to cacts.eth")
	bc.AddBlock("Send 2 more BTC to cacts.eth")

	bc.Print()

}
