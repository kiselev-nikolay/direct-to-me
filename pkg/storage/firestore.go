package storage

import (
	"context"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FireStoreStorage store key value data in google firestore
type FireStoreStorage struct {
	firestoreClient *firestore.Client
	caches          map[string]map[string]interface{}
	mutex           sync.Mutex
	cancelChannel   chan struct{}
}

type FireStoreStorageConf struct {
	ProjectID, CredentialsPath string
}

func (fs *FireStoreStorage) Connect(conf FireStoreStorageConf) error {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, conf.ProjectID, option.WithCredentialsFile(conf.CredentialsPath))
	if err != nil {
		return err
	}
	fs.caches = make(map[string]map[string]interface{})
	fs.firestoreClient = client
	return nil
}

func (fs *FireStoreStorage) ResetCache() {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	fs.caches = make(map[string]map[string]interface{})
}

func (fs *FireStoreStorage) Quit() {
	fs.firestoreClient.Close()
	close(fs.cancelChannel)
}

func (fs *FireStoreStorage) set(collection string, key string, value interface{}) error {
	docRef := fs.firestoreClient.Collection(collection).Doc(key)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	_, err := docRef.Set(ctx, value)
	if err != nil {
		return err
	}
	fs.mutex.Lock()
	col, ok := fs.caches[collection]
	if !ok {
		col = make(map[string]interface{})
		fs.caches[collection] = col
	}
	col[key] = value
	fs.mutex.Unlock()
	return nil
}

func (fs *FireStoreStorage) get(collection string, key string, value interface{}) (interface{}, error) {
	fs.mutex.Lock()
	col, ok := fs.caches[collection]
	if !ok {
		col = make(map[string]interface{})
		fs.caches[collection] = col
	}
	cachedValue, ok := col[key]
	fs.mutex.Unlock()
	if ok {
		return cachedValue, nil
	}
	docRef := fs.firestoreClient.Collection(collection).Doc(key)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	doc, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}
	err = doc.DataTo(value)
	if err != nil {
		return nil, err
	}
	fs.mutex.Lock()
	fs.caches[collection][key] = value
	fs.mutex.Unlock()
	return value, nil
}

func (fs *FireStoreStorage) delete(collection string, key string) error {
	docRef := fs.firestoreClient.Collection(collection).Doc(key)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	_, err := docRef.Delete(ctx)
	if err != nil {
		return err
	}
	fs.mutex.Lock()
	col, ok := fs.caches[collection]
	if !ok {
		col = make(map[string]interface{})
		fs.caches[collection] = col
	}
	delete(col, key)
	fs.mutex.Unlock()
	return nil
}

func (fs *FireStoreStorage) listKeys(collection string) (*firestore.DocumentIterator, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	iter := fs.firestoreClient.Collection(collection).Documents(ctx)
	return iter, cancel
}

func IsNotFoundError(err error) bool {
	return status.Code(err) == codes.NotFound
}

type Redirect struct {
	FromURI         string
	ToURL           string
	RedirectAfter   string
	URLTemplate     string
	MethodTemplate  string
	HeadersTemplate string
	BodyTemplate    string
}

func (fs *FireStoreStorage) SetRedirect(key string, value *Redirect) error {
	return fs.set("Redirect", key, value)
}

func (fs *FireStoreStorage) GetRedirect(key string) (*Redirect, error) {
	redirectPointer := &Redirect{}
	redirectValue, err := fs.get("Redirect", key, redirectPointer)
	if err != nil {
		return nil, err
	}
	redirectPointer = redirectValue.(*Redirect)
	return redirectPointer, err
}

func (fs *FireStoreStorage) DeleteRedirect(key string) error {
	return fs.delete("Redirect", key)
}

func (fs *FireStoreStorage) ListRedirects() ([]*Redirect, error) {
	redirects := make([]*Redirect, 0)
	documentIterator, cancel := fs.listKeys("Redirect")
	defer cancel()
	documents, err := documentIterator.GetAll()
	if err != nil {
		return nil, err
	}
	for _, document := range documents {
		redirect := &Redirect{}
		err := document.DataTo(redirect)
		if err != nil {
			log.Println("Broken node:", document.Ref.Path)
		}
		redirects = append(redirects, redirect)
	}
	return redirects, nil
}
