package filestore

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"reflect"
)

const defaultStorePath = "./data"

func Save[T any](ctx context.Context, name string, src T) error {
	collectionName := reflect.TypeOf(src).Name()
	dir := path.Join(defaultStorePath, collectionName)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("creating path: %w", err)
	}

	path := path.Join(dir, name+".json")

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating JSON file: %w", err)
	}

	defer file.Close()

	enc := json.NewEncoder(file)

	if err := enc.Encode(src); err != nil {
		return fmt.Errorf("enconding JSON: %w", err)
	}

	return nil
}

func RetrieveAll[T any](ctx context.Context) ([]T, error) {
	collectionName := reflect.TypeOf(*new(T)).Name()
	dir := path.Join(defaultStorePath, collectionName)

	fileNames, err := getFileNamesFromDir(dir)
	if err != nil {
		return nil, fmt.Errorf("retrieving file names: %w", err)
	}

	collection := make([]T, 0, len(fileNames))
	for _, name := range fileNames {
		err := func() error {
			path := path.Join(dir, name)
			file, err := os.Open(path)

			// Defer inside a loop is ok because we're using closure.
			defer file.Close()

			if err != nil {
				return fmt.Errorf("opening file: %w", err)
			}

			var item T
			dec := json.NewDecoder(file)
			if err := dec.Decode(&item); err != nil {
				return fmt.Errorf("decoding JSON: %w", err)
			}

			collection = append(collection, item)

			return nil
		}()

		if err != nil {
			return nil, err
		}
	}

	return collection, nil
}

// getFileNamesFromDir returns a slice of files that are located in the given directory.
func getFileNamesFromDir(dir string) ([]string, error) {
	var fileNames []string
	err := filepath.WalkDir(dir, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fileNames = append(fileNames, s)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileNames, nil
}
