package storage

type Storage interface {
	Get(string) string
	Set(string, string)
	Connect(interface{}) error
	ResetCache()
	Quit()
}
