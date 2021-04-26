package storage

type Storage interface {
	Set(collection string, key string, value interface{}) error
	Get(collection string, key string, value interface{}) error
	Delete(collection string, key string) error
	ListKeys(collection string) (v []string, err error)
	SetRedirect(key string, value *Redirect) error
	GetRedirect(key string) (*Redirect, error)
	DeleteRedirect(key string) error
	ListRedirects() ([]*Redirect, error)
}
