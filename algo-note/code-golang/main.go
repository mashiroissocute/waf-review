package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	sum := sha256.Sum256([]byte(""))
	fmt.Printf("%x\n", sum)
	fmt.Println(len(sum))
}
