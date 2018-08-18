package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Index        int       `json:"index"`
	Timestamp    time.Time `json:"timestamp"`
	Weight       int       `json:"weight"`
	Hash         string    `json:"hash"`
	PreviousHash string    `json:"previousHash"`
}

func main() {
	var blockchain []Block
	// create genesis block
	firstBlock := createFirstBlock()
	// append genesis block to blockchain to start it
	blockchain = append(blockchain, firstBlock)
	// let's get the json of the chain to see what it looks like
	jsonBlockChain, err := prettyBlockChain(blockchain)
	if err != nil {
		panic(err)
	}
	fmt.Println(jsonBlockChain)
	// get the first index of the first block
	i := 0
	// let's let a user enter values and update the chain with that value forever
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter a weight: ")
		weight, _ := reader.ReadString('\n')
		// convert the entered weight from string to int -- need to remove '\n'
		weightInt, err := strconv.Atoi(strings.Replace(weight, "\n", "", -1))
		if err != nil {
			panic(err)
		}
		// use the entered weight to create a block
		block, err := generateBlock(blockchain[i], weightInt)
		if err != nil {
			panic(err)
		}
		// add the new block to the chain
		blockchain = append(blockchain, *block)
		jsonBlockChain, err := prettyBlockChain(blockchain)
		if err != nil {
			panic(err)
		}
		// show the chain updated
		fmt.Println(jsonBlockChain)
		i++
	}
}

func createFirstBlock() Block {
	firstBlock := Block{
		Index:        0,
		Timestamp:    time.Now(),
		Weight:       0,
		PreviousHash: "",
	}
	firstBlockHash, err := calculateHash(firstBlock)
	if err != nil {
		panic(err)
	}
	firstBlock.Hash = firstBlockHash
	return firstBlock
}

func generateBlock(previousBlock Block, weight int) (*Block, error) {
	newBlock := Block{
		Index:        previousBlock.Index + 1,
		Timestamp:    time.Now(),
		Weight:       weight,
		PreviousHash: previousBlock.Hash,
	}
	hash, err := calculateHash(newBlock)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create a block for %d", weight)
	}
	newBlock.Hash = hash
	return &newBlock, nil
}

func calculateHash(block Block) (string, error) {
	// lets convert the block to JSON and use that for generating the hash
	bytes, err := json.Marshal(block)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal the block")
	}
	hash := sha256.New()
	_, err = hash.Write(bytes)
	if err != nil {
		return "", errors.Wrap(err, "failed to hash the block")
	}
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed), nil
}

func isBlockValid(currentBlock, previousBlock Block) (bool, error) {
	currentBlockHash, err := calculateHash(currentBlock)
	if err != nil {
		return false, errors.Wrap(err, "failed to compare blocks")
	}
	if (previousBlock.Index+1 != currentBlock.Index) ||
		(previousBlock.Hash != currentBlock.PreviousHash) ||
		(currentBlockHash != currentBlock.Hash) {
		return false, nil
	}
	return true, nil
}

func getLongestChain(block1, block2 []Block) []Block {
	if len(block1) > len(block2) {
		return block1
	}
	return block2
}

func prettyBlockChain(blockchain []Block) (string, error) {
	bytes, err := json.MarshalIndent(blockchain, "", "  ")
	if err != nil {
		return "", errors.Wrap(err, "failed to create json of blockchain")
	}
	return string(bytes), nil
}
