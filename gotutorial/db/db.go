package db

import (
	"fmt"
	"gotutorial/utils"
	"os"

	bolt "go.etcd.io/bbolt"
)

const (
	dbName            = "blockchain"
	dataBucket        = "data"
	blockBucket       = "block"
	transactionBucket = "tr"
	checkpoint        = "checkpoint"
)

var db *bolt.DB

func getDBName() string {
	return dbName + fmt.Sprint(os.Args[2][6:]) + ".db"
}
func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(getDBName(), 0600, nil)
		utils.HandleErr(err)
		db = dbPointer

		db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleErr(err)
			_, err = t.CreateBucketIfNotExists([]byte(blockBucket))
			utils.HandleErr(err)
			return err
		})
	}

	return db
}

func Close() {
	DB().Close()
}

func SaveBlock(hash string, data []byte) {
	fmt.Printf("Save Block: %s", data)
	err := DB().Update(func(t *bolt.Tx) error {
		b1 := t.Bucket([]byte(blockBucket))
		err := b1.Put([]byte(hash), data)
		return err
	})

	utils.HandleErr(err)
}

func SaveBlockchain(data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})
	utils.HandleErr(err)
}

func RestoreBlockchain() {}

func Checkpoint() []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
}

func Block(hash string) []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(blockBucket))
		data = b.Get([]byte(hash))
		return nil
	})
	return data
}

func EmptyBlocks() {
	DB().Update(func(t *bolt.Tx) error {
		t.Bucket([]byte(blockBucket))
		utils.HandleErr(t.DeleteBucket([]byte(blockBucket)))
		_, err := t.CreateBucket([]byte(blockBucket))
		utils.HandleErr(err)
		return nil
	})
}
