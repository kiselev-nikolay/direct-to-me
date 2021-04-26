package storage_test

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/kiselev-nikolay/direct-to-me/pkg/storage"
	"github.com/stretchr/testify/require"
)

var SavePath, _ = filepath.Abs(os.TempDir() + "/bitcaskdb")

type ExampleItem struct {
	Key      string
	Field    string
	Value    int
	IsActive bool
}

// ExampleItemCN collection name for first argument in storage.Storage Get, Set, etc
const ExampleItemCN = "ExampleItems"

func resetDB() {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	select {
	default:
		for {
			err = os.RemoveAll(SavePath)
			if err == nil {
				return
			}
			runtime.Gosched()
		}
	case <-ctx.Done():
		log.Fatal(err)
		return
	}
}

func TestReadWrite(t *testing.T) {
	require := require.New(t)
	strg := &storage.BitcaskStorage{
		SavePath: SavePath,
	}
	strg.Connect()
	defer func() {
		strg.Quit()
		resetDB()
	}()

	db := strg.GetDB()

	value := ExampleItem{Key: "key1", Field: "test", Value: 12, IsActive: true}
	err := strg.Set(ExampleItemCN, "key1", &value)
	require.NoError(err)
	require.True(db.Has([]byte(ExampleItemCN + ":" + "key1")))
	v, err := db.Get([]byte(ExampleItemCN + ":" + "key1"))
	require.NoError(err)
	require.NotEmpty(v)

	resultValue := ExampleItem{}
	err = strg.Get(ExampleItemCN, "key1", &resultValue)
	require.NoError(err)
	require.EqualValues(value, resultValue)
}

func TestIndexes(t *testing.T) {
	require := require.New(t)
	strg := &storage.BitcaskStorage{
		SavePath: SavePath,
	}
	strg.Connect()
	defer func() {
		strg.Quit()
		resetDB()
	}()

	db := strg.GetDB()

	value := ExampleItem{Key: "key1", Field: "test", Value: 12, IsActive: true}
	strg.Set(ExampleItemCN, "key1", &value)
	require.True(db.Has([]byte("index:" + ExampleItemCN)))
	index := make([]string, 0)
	err := strg.Get("index", ExampleItemCN, &index)
	require.NoError(err)
	require.Len(index, 1)
	require.Contains(index, "key1")

	err = strg.Delete(ExampleItemCN, "key1")
	require.NoError(err)
	index = make([]string, 0)
	err = strg.Get("index", ExampleItemCN, &index)
	require.NoError(err)
	require.Empty(index)
}

func TestIndexBug(t *testing.T) {
	// See: https://github.com/kiselev-nikolay/direct-to-me/issues/7
	require := require.New(t)
	strg := &storage.BitcaskStorage{
		SavePath: SavePath,
	}
	strg.Connect()
	defer func() {
		strg.Quit()
		resetDB()
	}()

	value := ExampleItem{Key: "key1", Field: "test", Value: 12, IsActive: true}
	strg.Set(ExampleItemCN, "key1", &value)
	strg.Set(ExampleItemCN, "key1", &value)
	v, err := strg.ListKeys(ExampleItemCN)
	require.NoError(err)
	require.Len(v, 1)
}
