package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/prologic/bitcask"
)

type BitcaskStorage struct {
	db *bitcask.Bitcask
}

var _ Storage = (*BitcaskStorage)(nil)

func (bs *BitcaskStorage) Connect() error {
	db, err := bitcask.Open("./bitcaskdb")
	if err != nil {
		return err
	}
	bs.db = db
	err = bs.db.Sync()
	return err
}
func (bs *BitcaskStorage) Quit() {
	bs.db.Close()
}
func (bs *BitcaskStorage) setNoIndex(collection string, key string, value interface{}) error {
	record := []byte{}
	encoder := gob.NewEncoder(bytes.NewBuffer(record))
	err := encoder.Encode(value)
	if err != nil {
		err = fmt.Errorf("bitcask set %v:%v->%v encoder error: %v", collection, key, value, err)
		return err
	}
	err = bs.db.Put([]byte(collection+":"+key), record)
	return err
}
func (bs *BitcaskStorage) set(collection string, key string, value interface{}) error {
	err := bs.addKey(collection, key)
	if err != nil {
		return err
	}
	err = bs.setNoIndex(collection, key, value)
	return err
}
func IsNotFoundError(err error) bool {
	return err == bitcask.ErrKeyNotFound
}
func (bs *BitcaskStorage) get(collection string, key string, value interface{}) error {
	record, err := bs.db.Get([]byte(collection + ":" + key))
	if err != nil {
		if err == bitcask.ErrKeyNotFound {
			return bitcask.ErrKeyNotFound
		}
		err = fmt.Errorf("bitcask get <-%v:%v read error: %v", collection, key, err)
		return err
	}
	log.Println(record)
	decoder := gob.NewDecoder(bytes.NewBuffer(record))
	err = decoder.Decode(value)
	if err != nil {
		err = fmt.Errorf("bitcask get <-%v:%v decode error: %v", collection, key, err)
		return err
	}
	return err
}
func (bs *BitcaskStorage) delete(collection string, key string) error {
	bs.deleteKey(collection, key)
	err := bs.db.Delete([]byte(collection + ":" + key))
	return err
}
func (bs *BitcaskStorage) listKeys(collection string) ([]string, error) {
	keys := make([]string, 0)
	err := bs.get("index", collection, &keys)
	return keys, err
}
func (bs *BitcaskStorage) addKey(collection string, key string) error {
	keys := make([]string, 0)
	err := bs.get("index", collection, &keys)
	if err != nil && !IsNotFoundError(err) {
		return err
	}
	keys = append(keys, key)
	err = bs.setNoIndex("index", collection, keys)
	return err
}
func (bs *BitcaskStorage) deleteKey(collection string, key string) error {
	keys := make([]string, 0)
	newKeys := make([]string, 0)
	err := bs.get("index", collection, &keys)
	if err != nil {
		return err
	}
	for _, newKey := range keys {
		if newKey != key {
			newKeys = append(newKeys, newKey)
		}
	}
	err = bs.setNoIndex("index", collection, &newKeys)
	return err
}
func (bs *BitcaskStorage) SetRedirect(key string, value *Redirect) error {
	err := bs.set("redirects", key, value)
	return err
}
func (bs *BitcaskStorage) GetRedirect(key string) (*Redirect, error) {
	redirect := &Redirect{}
	err := bs.get("redirects", key, redirect)
	if err != nil {
		return nil, err
	}
	return redirect, nil
}
func (bs *BitcaskStorage) DeleteRedirect(key string) error {
	err := bs.delete("redirects", key)
	return err
}
func (bs *BitcaskStorage) ListRedirects() ([]*Redirect, error) {
	redirects := make([]*Redirect, 0)
	redirectKeys, err := bs.listKeys("redirects")
	if err != nil {
		return nil, err
	}
	for _, key := range redirectKeys {
		redirect := &Redirect{}
		bs.get("redirects", key, redirect)
		redirects = append(redirects, redirect)
	}
	return redirects, nil
}
