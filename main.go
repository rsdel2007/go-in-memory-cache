package main

import (
	"fmt"
	"log"

	"github.com/rsdel2007/go-in-memory-cache/cache"
)

func main() {
	db := cache.NewSimpleDB()
	fmt.Println("Example 1:")

	db.Begin()      // start transaction 1
	db.Set("a", 10) // set a 10
	value, exists := db.Get("a")
	fmt.Printf("Get a: %v (exists : %v)\n", value, exists) // Expected: 10

	db.Begin()
	db.Set("a", 20)
	value, exists = db.Get("a")
	fmt.Printf("GET a: %v (exists: %v)\n", value, exists) // Expected: 20

	err := db.Rollback() // Rollback transaction 2
	if err != nil {
		log.Println(err)
	}

	value, exists = db.Get("a")
	fmt.Printf("GET a after rollback: %v (exists: %v)\n", value, exists) // Expected: 10

	err = db.Rollback() // Rollback transaction 1
	if err != nil {
		log.Println(err)
	}
	value, exists = db.Get("a")
	fmt.Printf("GET a after second rollback: %v (exists: %v)\n", value, exists) // Expected: NULL
	fmt.Println("Example 2:")
	db.Begin()      // Start transaction 1
	db.Set("a", 30) // Set a = 30
	db.Begin()      // Start transaction 2 (nested)
	db.Set("a", 40) // Set a = 40 in transaction 2

	err = db.Commit() // Commit
	if err != nil {
		log.Println(err)
	}
	value, exists = db.Get("a")
	fmt.Printf("GET a after commit: %v (exists: %v)\n", value, exists) // Expected: 40

	err = db.Rollback() // Attempt rollback after commit
	if err != nil {
		fmt.Println("ROLLBACK:", err) // Expected: "NO TRANSACTION" error
	}

	fmt.Println("\nExample 3:")
	db.Set("a", 50) // Set a = 50 outside transaction
	db.Begin()      // Start transaction 1
	value, exists = db.Get("a")
	fmt.Printf("GET a: %v (exists: %v)\n", value, exists) // Expected: 50

	db.Set("a", 60) // Set a = 60 within transaction 1
	db.Begin()      // Start transaction 2 (nested)
	db.Unset("a")   // Unset a within transaction 2
	value, exists = db.Get("a")
	fmt.Printf("GET a after unset: %v (exists: %v)\n", value, exists) // Expected: NULL

	err = db.Rollback() // Rollback transaction 2
	if err != nil {
		log.Println(err)
	}
	value, exists = db.Get("a")
	fmt.Printf("GET a after rollback: %v (exists: %v)\n", value, exists) // Expected: 60

	err = db.Commit() // Commit transaction 1
	if err != nil {
		log.Println(err)
	}
	value, exists = db.Get("a")
	fmt.Printf("GET a after final commit: %v (exists: %v)\n", value, exists) // Expected: 60

}

// where is grpc?
// config -> application level
// cluster1 -> cluster2
