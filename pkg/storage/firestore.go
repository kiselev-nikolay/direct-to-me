package storage

import (
	"context"
	"sync"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

// FireStoreStorage store key value data in google firestore
type FireStoreStorage struct {
	firestoreClient *firestore.Client
	caches          FireStoreStorageCaches
	mutex           sync.Mutex
	cancelChannel   chan struct{}
}

type FireStoreStorageCaches struct {
	fl map[string]string
}
type FireStoreStorageConf struct {
	projectID, credentialsPath string
}

func (kv *FireStoreStorage) Connect(conf FireStoreStorageConf) error {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, conf.projectID, option.WithCredentialsFile(conf.credentialsPath))
	if err != nil {
		return err
	}
	kv.caches.fl = make(map[string]string)
	kv.firestoreClient = client
	return nil
}

func (kv *FireStoreStorage) Set(k, v string) {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	kv.caches.fl[k] = v
	// Todo send to cloud
}

func (kv *FireStoreStorage) Get(k string) string {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	v := kv.caches.fl[k]
	// Todo get from cloud
	return v
}

func (kv *FireStoreStorage) ResetCache() {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	kv.caches.fl = make(map[string]string)
}
