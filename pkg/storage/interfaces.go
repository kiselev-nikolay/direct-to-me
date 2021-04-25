package storage

type Storage interface {
	set(collection string, key string, value interface{}) error
	get(collection string, key string, value interface{}) error
	delete(collection string, key string) error
	listKeys(collection string) (v []string, err error)
	SetRedirect(key string, value *Redirect) error
	GetRedirect(key string) (*Redirect, error)
	DeleteRedirect(key string) error
	ListRedirects() ([]*Redirect, error)
}
