package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Block struct { // Define a block struct
	data         map[string]interface{} // creates key value, key: string, value: interface (any type)
	hash         string
	previousHash string
	timestamp    time.Time // gives the current time
	pow          int
}

type Blockchain struct { // Define a blockchain struct
	genesisBlock Block
	chain        []Block
	difficulty   int
}

func (b Block) calculateHash() string { // method to calculate the hash of a block, read-only -> takes a copy of block
	data, _ := json.Marshal(b.data)                                                         // Convert data to JSON
	blockData := b.previousHash + string(data) + b.timestamp.String() + strconv.Itoa(b.pow) // Concatenate previousHash, data, timestamp and pow // itoa converts int to string
	blockHash := sha256.Sum256([]byte(blockData))                                           // Hash the concatenated string w/ sha256 algorithm
	return fmt.Sprintf("%x", blockHash)                                                     // Return the base 16 hash as a string
}

func (b *Block) mine(difficulty int) { // increments the pow and calc block hash until valid (# of 0s > difficulty), modifies the block -> takes a reference to the block, *Block is a receiver type
	for !strings.HasPrefix(b.hash, strings.Repeat("0", difficulty)) { // While the hash doesn't start with the required number of 0s
		b.pow++                    // Increment the pow
		b.hash = b.calculateHash() // Calculate the hash of the block
	}
}

func CreateBlockchain(difficulty int) Blockchain { // Create a new blockchain with a genesis block
	genesisBlock := Block{ // create instance of block -> assign to genesisBlock variable
		hash:      "0", // set hash of genesis block to 0, first block -> no previous hash, data is empty
		timestamp: time.Now(),
	}
	return Blockchain{ // create a new instance of the blockchain
		genesisBlock,          // store the genesis block as a field
		[]Block{genesisBlock}, // add genesis block to the chain (slice of blocks)
		difficulty,
	}
}

func (b *Blockchain) addBlock(from, to string, amount float64) { // *Blockchain receiver type, from, to, amount -> var parameters
	blockData := map[string]interface{}{ // Create blockData var (map) -> assign keys and values
		"from":   from,
		"to":     to,
		"amount": amount,
	}
	lastBlock := b.chain[len(b.chain)-1] // Get the last block on the chain
	newBlock := Block{
		data:         blockData,
		previousHash: lastBlock.hash,
		timestamp:    time.Now(),
	}
	newBlock.mine(b.difficulty)         // Mine the new block
	b.chain = append(b.chain, newBlock) // Append the new block to the chain
}

func (b Blockchain) isValid() bool {
	for i := range b.chain[1:] { // Loops through chain starting from the 2nd block, assign i the indexes which represent the blocks starting from the 2nd block.
		previousBlock := b.chain[i]                                                                               // get previous block starting from the genesis block
		currentBlock := b.chain[i+1]                                                                              // get current block starting from the 2nd block
		if currentBlock.hash != currentBlock.calculateHash() || currentBlock.previousHash != previousBlock.hash { // check validation
			return false
		}
	}
	return true
}

func main() {
	blockchain := CreateBlockchain(2) // create a new blockchain with a difficulty of 2
	blockchain.addBlock("Bob", "Alice", 5)
	blockchain.addBlock("John", "Bob", 15)

	fmt.Println(blockchain.isValid())

	fmt.Println("\nBlockchain Structure:")
	for i, block := range blockchain.chain {
		fmt.Printf("\nBlock %d:\n", i)
		fmt.Printf("  Previous Hash: %s\n", block.previousHash)
		fmt.Printf("  Data: %+v\n", block.data)
		fmt.Printf("  Timestamp: %s\n", block.timestamp)
		fmt.Printf("  PoW: %d\n", block.pow)
		fmt.Printf("  Hash: %s\n", block.hash)
	}
}
