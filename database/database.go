package database

import (
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
)

const HyperclairDB = "hyperclair.db"
const RegistryBucket = "Registries"

func InsertRegistryMapping(layerDigest string, registryURI string) error {
	db, err := open(HyperclairDB)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(RegistryBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		fmt.Printf("Saving %s[%s]\n", layerDigest, registryURI)
		err = b.Put([]byte(layerDigest), []byte(registryURI))

		return err
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
		return "", err
	}
	if value == nil {
		return "", errors.New(layerDigest + " Mapping not found")
	}
	return string(value), nil
}

func open(dbName string) (*bolt.DB, error) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, err
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(RegistryBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return err
	})
	return db, nil
}

func IsHealthy() interface{} {
	type Health struct {
		IsHealthy bool
	}

	db, err := open(HyperclairDB)
	defer db.Close()

	if err != nil {
		return Health{false}
	}

	return Health{true}
}
