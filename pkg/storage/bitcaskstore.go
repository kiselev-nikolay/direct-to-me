package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/prologic/bitcask"
)

type BitcaskStorage struct {
	db       *bitcask.Bitcask
	SavePath string
}

var _ Storage = (*BitcaskStorage)(nil)

func (bs *BitcaskStorage) GetDB() *bitcask.Bitcask {
	return bs.db
}
func (bs *BitcaskStorage) Connect() error {
	if bs.SavePath == "" {
		bs.SavePath = "./bitcaskdb"
	}
	db, err := bitcask.Open(bs.SavePath)
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
	record := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(record)
	err := encoder.Encode(value)
	if err != nil {
		err = fmt.Errorf("bitcask set %v:%v->%v encoder error: %v", collection, key, value, err)
		return err
	}
	err = bs.db.Put([]byte(collection+":"+key), record.Bytes())
	return err
}
func (bs *BitcaskStorage) Set(collection string, key string, value interface{}) error {
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
func (bs *BitcaskStorage) Get(collection string, key string, value interface{}) error {
	record, err := bs.db.Get([]byte(collection + ":" + key))
	if err != nil {
		if err == bitcask.ErrKeyNotFound {
			return bitcask.ErrKeyNotFound
		}
		err = fmt.Errorf("bitcask get <-%v:%v read error: %v", collection, key, err)
		return err
	}
	decoder := gob.NewDecoder(bytes.NewBuffer(record))
	err = decoder.Decode(value)
	if err != nil {
		err = fmt.Errorf("bitcask get <-%v:%v decode error: %v", collection, key, err)
		return err
	}
	return err
}
func (bs *BitcaskStorage) Delete(collection string, key string) error {
	bs.deleteKey(collection, key)
	err := bs.db.Delete([]byte(collection + ":" + key))
	return err
}
func (bs *BitcaskStorage) ListKeys(collection string) ([]string, error) {
	keys := make([]string, 0)
	err := bs.Get("index", collection, &keys)
	return keys, err
}
func (bs *BitcaskStorage) addKey(collection string, key string) error {
	keys := make([]string, 0)
	err := bs.Get("index", collection, &keys)
	if err != nil && !IsNotFoundError(err) {
		return err
	}
	for _, i := range keys {
		if i == key {
			return nil
		}
	}
	keys = append(keys, key)
	err = bs.setNoIndex("index", collection, keys)
	return err
}
func (bs *BitcaskStorage) deleteKey(collection string, key string) error {
	keys := make([]string, 0)
	newKeys := make([]string, 0)
	err := bs.Get("index", collection, &keys)
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
	err := bs.Set("redirects", key, value)
	return err
}
func (bs *BitcaskStorage) GetRedirect(key string) (*Redirect, error) {
	redirect := &Redirect{}
	err := bs.Get("redirects", key, redirect)
	if err != nil {
		return nil, err
	}
	return redirect, nil
}
func (bs *BitcaskStorage) DeleteRedirect(key string) error {
	err := bs.Delete("redirects", key)
	return err
}
func (bs *BitcaskStorage) ListRedirects() ([]*Redirect, error) {
	redirects := make([]*Redirect, 0)
	redirectKeys, err := bs.ListKeys("redirects")
	if err != nil {
		return nil, err
	}
	for _, key := range redirectKeys {
		redirect := &Redirect{}
		bs.Get("redirects", key, redirect)
		redirects = append(redirects, redirect)
	}
	return redirects, nil
}
