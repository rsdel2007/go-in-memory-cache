package cache

import "testing"

func TestSetWithoutTransaction(t *testing.T) {
	db := NewSimpleDB()

	// Perform a set operation without any transaction
	db.Set("key1", 10)
	value, exists := db.Get("key1")
	if !exists || value != 10 {
		t.Errorf("Expected key1 to be 10, got %d", value)
	}
}

func TestUnsetWithoutTransaction(t *testing.T) {
	db := NewSimpleDB()

	db.Set("key1", 10)
	db.Unset("key1")

	_, exists := db.Get("key1")
	if exists {
		t.Errorf("Expected key1 to be unset")
	}
}

func TestCommitWithNestedTransactions(t *testing.T) {
	db := NewSimpleDB()

	db.Set("key1", 10)

	db.Begin() // Transaction 1
	db.Set("key1", 20)

	db.Begin() // Nested Transaction 2
	db.Set("key1", 30)

	err := db.Commit() // Commit both transactions
	if err != nil {
		t.Errorf("Unexpected error on commit: %v", err)
	}

	value, exists := db.Get("key1")
	if !exists || value != 30 {
		t.Errorf("Expected key1 to be 30 after commit, got %d", value)
	}
}

func TestRollbackWithMultipleOperations(t *testing.T) {
	db := NewSimpleDB()

	db.Set("key1", 10)
	db.Begin()

	db.Set("key1", 20)
	db.Set("key2", 30)

	err := db.Rollback()
	if err != nil {
		t.Errorf("Unexpected error on rollback: %v", err)
	}

	// After rollback, `key1` should still be 10, and `key2` should not exist
	value, exists := db.Get("key1")
	if !exists || value != 10 {
		t.Errorf("Expected key1 to be 10 after rollback, got %d", value)
	}

	_, exists = db.Get("key2")
	if exists {
		t.Errorf("Expected key2 to not exist after rollback")
	}
}

func TestRollbackAfterCommit(t *testing.T) {
	db := NewSimpleDB()

	db.Set("key1", 10)
	db.Begin()

	db.Set("key1", 20)

	err := db.Commit()
	if err != nil {
		t.Errorf("Unexpected error on commit: %v", err)
	}

	// Try to rollback after commit
	err = db.Rollback()
	if err == nil || err.Error() != "NO TRANSACTION" {
		t.Errorf("Expected NO TRANSACTION error after commit, got %v", err)
	}

	value, exists := db.Get("key1")
	if !exists || value != 20 {
		t.Errorf("Expected key1 to be 20 after commit, got %d", value)
	}
}

func TestNestedTransactionCommitAndRollback(t *testing.T) {
	db := NewSimpleDB()

	db.Set("key1", 10)
	db.Begin() // Start outer transaction

	db.Set("key1", 20)
	db.Begin() // Start nested transaction

	db.Set("key1", 30)
	err := db.Commit() // Commit nested transaction

	if err != nil {
		t.Errorf("Unexpected error on nested commit: %v", err)
	}

	value, exists := db.Get("key1")
	if !exists || value != 30 {
		t.Errorf("Expected key1 to be 30 after nested commit, got %d", value)
	}

	err = db.Rollback() // Rollback outer transaction
	if err == nil {
		t.Errorf("Expected err but didn't get one on rollback: %v", err)
	}

	value, exists = db.Get("key1")
	if !exists || value != 30 {
		t.Errorf("Expected key1 to remain 30 after outer rollback, got %d", value)
	}
}

func TestTransactionAfterMultipleCommits(t *testing.T) {
	db := NewSimpleDB()

	db.Begin()
	db.Set("key1", 10)
	err := db.Commit()
	if err != nil {
		t.Errorf("Unexpected error on first commit: %v", err)
	}

	db.Begin()
	db.Set("key1", 20)
	err = db.Commit()
	if err != nil {
		t.Errorf("Unexpected error on second commit: %v", err)
	}

	value, exists := db.Get("key1")
	if !exists || value != 20 {
		t.Errorf("Expected key1 to be 20 after multiple commits, got %d", value)
	}
}

func TestRollbackNestedTransactionWithUnset(t *testing.T) {
	db := NewSimpleDB()

	db.Set("key1", 10)
	db.Begin()

	db.Set("key1", 20)
	db.Begin()
	db.Unset("key1")

	_, exists := db.Get("key1")
	if exists {
		t.Errorf("Expected key1 to be unset")
	}

	err := db.Rollback()
	if err != nil {
		t.Errorf("Unexpected error on rollback: %v", err)
	}

	// After rollback, key1 should still be 20
	value, exists := db.Get("key1")
	if !exists || value != 20 {
		t.Errorf("Expected key1 to be 20 after rollback, got %d", value)
	}

	err = db.Commit()
	if err != nil {
		t.Errorf("Unexpected error on commit: %v", err)
	}

	value, exists = db.Get("key1")
	if !exists || value != 20 {
		t.Errorf("Expected key1 to be 20 after commit, got %d", value)
	}
}

func TestCommitOnEmptyTransactionStack(t *testing.T) {
	db := NewSimpleDB()

	err := db.Commit()
	if err == nil || err.Error() != "NO TRANSACTION" {
		t.Errorf("Expected NO TRANSACTION error when committing with empty stack, got %v", err)
	}
}

func TestRollbackOnEmptyTransactionStack(t *testing.T) {
	db := NewSimpleDB()

	err := db.Rollback()
	if err == nil || err.Error() != "NO TRANSACTION" {
		t.Errorf("Expected NO TRANSACTION error when rolling back with empty stack, got %v", err)
	}
}

func TestSetAndGetMultipleKeys(t *testing.T) {
	db := NewSimpleDB()

	db.Set("key1", 10)
	db.Set("key2", 20)
	db.Set("key3", 30)

	value1, exists1 := db.Get("key1")
	if !exists1 || value1 != 10 {
		t.Errorf("Expected key1 to be 10, got %d", value1)
	}

	value2, exists2 := db.Get("key2")
	if !exists2 || value2 != 20 {
		t.Errorf("Expected key2 to be 20, got %d", value2)
	}

	value3, exists3 := db.Get("key3")
	if !exists3 || value3 != 30 {
		t.Errorf("Expected key3 to be 30, got %d", value3)
	}
}
