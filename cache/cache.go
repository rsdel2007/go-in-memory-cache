package cache

import (
	"errors"
)

type SimpleDb interface {
	Set(key string, value int)
	Get(key string) (int, bool)
	Unset(key string)
	Begin()
	Commit() error
	Rollback() error
}

type simpleDB struct {
	data         map[string]int
	transactions []map[string]*int
}

func (db *simpleDB) Set(key string, value int) {
	if len(db.transactions) > 0 {
		txn := db.transactions[len(db.transactions)-1]
		if _, exists := db.data[key]; !exists {
			// key does not exist before, store nil
			txn[key] = nil
		} else if _, exists := txn[key]; !exists {
			v := db.data[key]
			txn[key] = &v
		}
	}
	db.data[key] = value
}

func (db *simpleDB) Get(key string) (int, bool) {
	value, exists := db.data[key]
	return value, exists
}

func (db *simpleDB) Unset(key string) {
	if len(db.transactions) > 0 {
		txn := db.transactions[len(db.transactions)-1]
		if _, exists := db.data[key]; exists && txn[key] == nil {
			v := db.data[key]
			txn[key] = &v
		} else {
			txn[key] = nil
		}
	}
	delete(db.data, key)
}

func (db *simpleDB) Begin() {
	db.transactions = append(db.transactions, make(map[string]*int))
}

func (db *simpleDB) Commit() error {
	if len(db.transactions) == 0 {
		return errors.New("NO TRANSACTION")
	}
	// Clear all transactions after committing
	db.transactions = []map[string]*int{}
	return nil
}

func (db *simpleDB) Rollback() error {
	if len(db.transactions) == 0 {
		return errors.New("NO TRANSACTION")
	}

	lastTxn := db.transactions[len(db.transactions)-1]
	db.transactions = db.transactions[:len(db.transactions)-1]

	for key, originalValue := range lastTxn {
		if originalValue == nil {
			// if orignal was nil, delete key
			delete(db.data, key)
		} else {
			db.data[key] = *originalValue
		}
	}
	return nil
}

func NewSimpleDB() SimpleDb {
	return &simpleDB{
		data:         make(map[string]int),
		transactions: []map[string]*int{},
	}
}
