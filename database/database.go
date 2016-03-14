package database

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const HyperclairDB = "/data/hyperclair.db"
const RegistryBucket = "Registries"

func InsertRegistryMapping(layerDigest string, registryURI string) error {
	db, err := open(HyperclairDB)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		log.Printf("Saving %s[%s]\n", layerDigest, registryURI)
		err = tx.Bucket([]byte(RegistryBucket)).Put([]byte(layerDigest), []byte(registryURI))
		if err != nil {
			return fmt.Errorf("adding registry mapping: %v", err)
		}
		return nil
	})

}
func List() {
	db, err := open("hyperclair.db")
	defer db.Close()
	if err != nil {
		log.Fatalf("db list: %v", err)
	}
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(RegistryBucket))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}

		return nil
	})
}

func GetRegistryMapping(layerDigest string) (string, error) {
	db, err := open(HyperclairDB)
	defer db.Close()

	if err != nil {
		return "", err
	}

	var value []byte

	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(RegistryBucket))
		value = b.Get([]byte(layerDigest))

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("retrieving registry mapping: %v", err)
	}
	if value == nil {
		return "", fmt.Errorf("%v mapping not found", layerDigest)
	}
	return string(value), nil
}

func open(dbName string) (*bolt.DB, error) {
	db, err := bolt.Open(dbName, 0600, nil)

	if err != nil {
		return nil, fmt.Errorf("opening db: %v", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(RegistryBucket))
		if err != nil {
			return fmt.Errorf("creating bucket: %v", err)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("updating db: %v", err)
	}
	return db, nil
}

func IsHealthy() (interface{}, bool) {
	type Health struct {
		IsHealthy bool
	}

	db, err := open(HyperclairDB)
	if err != nil {
		return Health{false}, false
	}

	defer db.Close()

	return Health{true}, true
}
