/*
Package filestore offers a concurrency-safe generic store for any item type.
Items are stored as JSON files in a
specified directory, with each item type being stored in its own subdirectory.
The package supports basic operations
like storing a new item and fetching all stored items.

Instances of FileStore are safe for concurrent use, achieved by using a mutex lock whenever
accessing the file system.
The name of the JSON file used to store an item is provided when the item is stored.
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
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sync"
)

var (
	// ErrFileExists is the error returned when trying
	// to store an item with a name that already exists.
	ErrFileExists = fs.ErrExist

	// ErrNotExist is the error returned
	// when trying to fetch an item that does not exist.
	ErrNotExist = os.ErrNotExist
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
func (f *FileStore[T]) Store(name string, item T) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if err := os.MkdirAll(f.dir, os.ModePerm); err != nil {
		return fmt.Errorf("creating path: %w", err)
	}

	pth := path.Join(f.dir, name)
	if info, err := os.Stat(pth); !os.IsNotExist(err) {
		log.Println(info)
		return ErrFileExists
	}

	file, err := os.Create(pth)
	if err != nil {
		return fmt.Errorf("creating JSON file: %w", err)
	}

	defer file.Close()

	if err := json.NewEncoder(file).Encode(item); err != nil {
		if err := os.Remove(pth); err != nil {
			return fmt.Errorf("removing JSON file: %w", err)
		}

		return fmt.Errorf("encoding JSON: %w", err)
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
		return nil, fmt.Errorf("%w: %w", ErrNotExist, err)
	}

	return coll, nil
}
