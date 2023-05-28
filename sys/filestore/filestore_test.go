package filestore

import (
	"fmt"
	"os"
	"path"
	"testing"
)

type Email struct {
	Address string
}

func TestFileStore(t *testing.T) {
	tt := []struct {
		name         string
		testFunction func(t *testing.T)
	}{
		{name: "TestNew", testFunction: TestNew},
		{name: "TestStore", testFunction: TestStore},
		{name: "TestStoreFileExists", testFunction: TestStoreFileExists},
		{name: "TestFetchAll", testFunction: TestFetchAll},
	}

	for _, tc := range tt {
		t.Run(tc.name, tc.testFunction)
	}
}

func TestNew(t *testing.T) {
	fs := New[Email]("./tmp/test_path")
	defer func() {
		if err := os.RemoveAll("./tmp/test_path"); err != nil {
			t.Errorf("error while removing test_path: %v", err)
		}
	}()

	expectedDir := path.Join("./tmp/test_path", "Email")
	if fs.dir != expectedDir {
		t.Errorf("expected dir to be %s, got %s", expectedDir, fs.dir)
	}
}

func TestStore(t *testing.T) {
	fs := New[Email]("./tmp/test_path")
	defer func() {
		if err := os.RemoveAll("./tmp/test_path"); err != nil {
			t.Errorf("error while removing test_path: %v", err)
		}
	}()

	item := Email{Address: "test@example.com"}

	err := fs.Store("item1", item)
	if err != nil {
		t.Errorf("Error while storing item: %v", err)
	}

	if _, err := os.Stat(path.Join(fs.dir, "item1")); os.IsNotExist(err) {
		t.Errorf("File was not created")
	}
}

func TestStoreFileExists(t *testing.T) {
	fs := New[Email]("./tmp/test_path")
	defer func() {
		if err := os.RemoveAll("./tmp/test_path"); err != nil {
			t.Errorf("error while removing test_path: %v", err)
		}
	}()

	item := Email{Address: "test@example.com"}
	fs.Store("item1", item)

	err := fs.Store("item1", item)
	if err != ErrFileExists {
		t.Errorf("Expected ErrFileExists, got %v", err)
	}
}

func TestFetchAll(t *testing.T) {
	fs := New[Email]("./tmp/test_path")
	defer func() {
		if err := os.RemoveAll("./tmp/test_path"); err != nil {
			t.Errorf("error while removing test_path: %v", err)
		}
	}()

	items := []Email{
		{Address: "test1@example.com"},
		{Address: "test2@example.com"},
		{Address: "test3@example.com"},
	}

	for i, item := range items {
		err := fs.Store(fmt.Sprintf("item%d", i+1), item)
		if err != nil {
			t.Errorf("Error while storing item: %v", err)
		}
	}

	fetchedItems, err := fs.FetchAll()
	if err != nil {
		t.Errorf("Error while fetching all items: %v", err)
	}

	if len(fetchedItems) != len(items) {
		t.Errorf("Expected %d items, got %d", len(items), len(fetchedItems))
	}

	for i, item := range fetchedItems {
		if item.Address != items[i].Address {
			t.Errorf("Expected item at index %d to be %s, got %s", i, items[i].Address, item.Address)
		}
	}
}
