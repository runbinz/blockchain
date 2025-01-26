package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Block struct { // Create a block struct
	data         map[string]interface{} // creates key value, key: string, value: interface (any type)
	hash         string
	previousHash string
	timestamp    time.Time // gives the current time
	pow          int
}

type Blockchain struct { // Create a blockchain struct
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

func (b *Block) mine(difficulty int) { // increments the pow and calc block hash until valid (# of 0s > difficulty), modifies the block -> takes a reference to the block
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
