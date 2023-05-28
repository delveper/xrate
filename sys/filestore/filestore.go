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

var ErrFileExists = fs.ErrExist

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
		return fmt.Errorf("enconding JSON: %w", err)
	}

	return nil
}

func (f *FileStore[T]) FetchAll() ([]T, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	var coll []T

	walkFn := func(pth string, ent fs.DirEntry, err error) error {
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

	if err := filepath.WalkDir(f.dir, walkFn); err != nil {
		return nil, err
	}

	return coll, nil
}
