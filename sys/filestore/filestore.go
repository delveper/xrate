/*
Package filestore offers a concurrency-safe generic store for any item type.
Items are stored as JSON files in a
specified directory, with each item type being stored in its own subdirectory.
The package supports basic operations
like storing a new item and fetching all stored items.

Instances of FileStore are safe for concurrent use, achieved by using a mutex lock whenever
accessing the file system.
The name of the JSON file is hashed with the name of the item.
If a file with the same name already exists, an error is returned.

Example usage:

	type Person struct {
		Name string
		Age  int
	}

	store := filestore.New[Person]("/path/to/store")
	err := store.Store("johndoe", Person{"John Doe", 30})

	persons, err := store.FetchAll()

The above would create a JSON file at "/path/to/store/Person/johndoe", containing the JSON representation of the
specified Person struct.
It could then retrieve all Person structs stored in the "/path/to/store/Person" directory.
*/
package filestore

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
)

type FileStore[T any] struct {
	mu  sync.Mutex
	dir string
}

func New[T any](pth string) *FileStore[T] {
	name := reflect.TypeOf(*new(T)).Name()
	dir := path.Join(pth, name)

	return &FileStore[T]{
		dir: dir,
	}
}

// Store method stores the item of type T as a JSON file with the given name.
func (f *FileStore[T]) Store(item T) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if reflect.ValueOf(item).IsZero() || f.dir == "" {
		return os.ErrInvalid
	}

	if err := os.MkdirAll(f.dir, os.ModePerm); err != nil {
		return fmt.Errorf("creating path: %w", err)
	}

	pth := path.Join(f.dir, hash(item))
	if _, err := os.Stat(pth); !os.IsNotExist(err) {
		return os.ErrExist
	}

	file, err := os.Create(pth)
	if err != nil {
		return fmt.Errorf("creating JSON file: %w", err)
	}

	defer file.Close()

	if err := json.NewEncoder(file).Encode(item); err != nil {
		err = fmt.Errorf("encoding JSON: %w", err)

		if errRm := os.Remove(pth); errRm != nil {
			return fmt.Errorf("%w: removing JSON file: %w", err, errRm)
		}

		return err
	}

	return nil
}

// FetchAll method fetches all items of type T stored in the FileStore as a slice.
func (f *FileStore[T]) FetchAll() ([]T, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	var coll []T

	// walkFn is the function applied to every file in the directory.
	walkDirFunc := func(pth string, ent fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walking path: %w", err)
		}

		if !ent.IsDir() {
			file, err := os.Open(pth)
			if err != nil {
				return fmt.Errorf("opening file: %w", err)
			}
			defer file.Close()

			var item T
			if err := json.NewDecoder(file).Decode(&item); err != nil {
				return fmt.Errorf("decoding JSON: %w", err)
			}

			coll = append(coll, item)
		}

		return nil
	}

	if err := filepath.WalkDir(f.dir, walkDirFunc); err != nil {
		return nil, fmt.Errorf("%w: %w", os.ErrNotExist, err)
	}

	return coll, nil
}

func (f *FileStore[T]) StoreAll(items ...T) error {
	for _, item := range items {
		if err := f.Store(item); err != nil {
			return fmt.Errorf("storing items: %w", err)
		}
	}

	return nil
}

func hash(item any) string {
	h := hmac.New(sha256.New, []byte("my_secret"))
	h.Write(fmt.Append(nil, item))

	return hex.EncodeToString(h.Sum(nil))
}

func TestSetup[T any](t *testing.T) (*FileStore[T], func()) {
	t.Helper()

	dir := t.TempDir()
	store := New[T](dir)

	teardown := func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Errorf("setup: removing temp dir: %v", err)
		}
	}

	return store, teardown
}
