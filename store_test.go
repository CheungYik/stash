package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathTransformFunc(t *testing.T) {
	key := "hello"
	got := CASPathTransformFunc(key)
	expectedOriginal := "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"
	expectedPathName := "aaf4c/61ddc/c5e8a/2dabe/de0f3/b482c/d9aea/9434d"
	assert.Equal(t, expectedOriginal, got.FileName)
	assert.Equal(t, expectedPathName, got.PathName)
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		Root:          "./data",
		PathTransform: CASPathTransformFunc,
	}
	store := NewStore(opts)
	key := "hello"
	content := "some data"

	reader := bytes.NewReader([]byte(content))
	if err := store.writeStream(key, reader); err != nil {
		t.Error(err)
	}

	if ok := store.Has(key); !ok {
		t.Errorf("expected key %s to exist", key)
	}

	read, err := store.Read(key)
	if err != nil {
		t.Error(err)
	}

	data, err := io.ReadAll(read)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, content, string(data))

	//err = store.Delete(key)
	//if err != nil {
	//	t.Error(err)
	//}
}
