package main

import (
	"fmt"
	"github.com/ahmedennaifer/taskq/pkg"
)

func main() {
	workerID := "fcbe46b6-a666-485b-b482-6fe1e2c459b9"
	secret := "test"
	
	hash1, _ := pkg.Hash(workerID, secret)
	fmt.Println("Hash 1:", hash1)
	
	hash2, _ := pkg.Hash(workerID, secret)
	fmt.Println("Hash 2:", hash2)
	
	fmt.Println("Match:", hash1 == hash2)
}
