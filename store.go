package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
)

const defaultRootPath = "/data/stash"

type PathKey struct {
	PathName string
	FileName string
}

// FirstPathName returns the first path name of the PathKey
func (slf PathKey) FirstPathName() string {
	paths := strings.Split(slf.PathName, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

// FullPath returns the full path of the PathKey
func (slf PathKey) FullPath() string {
	return path.Join(slf.PathName, slf.FileName)
}

// PathTransformFunc is a function that transforms a key into a PathKey
type PathTransformFunc func(string) PathKey

// DefaultPathTransformFunc is the default PathTransformFunc
var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		FileName: key,
	}
}

// CASPathTransformFunc is a PathTransformFunc that transforms a key into a PathKey using CAS(Content Addressable Storage) algorithm
var CASPathTransformFunc = func(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i+1)*blockSize
		paths[i] = hashStr[from:to]
	}
	return PathKey{
		PathName: strings.Join(paths, "/"),
		FileName: hashStr,
	}
}

type StoreOpts struct {
	Root          string
	PathTransform PathTransformFunc
}

type Store struct {
	StoreOpts
}

// NewStore creates a new Store
func NewStore(opts StoreOpts) *Store {
	if opts.PathTransform == nil {
		opts.PathTransform = DefaultPathTransformFunc
	}
	if len(opts.Root) == 0 {
		opts.Root = defaultRootPath
	}
	return &Store{
		StoreOpts: opts,
	}
}

// Has checks if a key exists in the store
func (slf *Store) Has(key string) bool {
	pathKey := slf.PathTransform(key)
	name := path.Join(slf.Root, pathKey.FullPath())
	_, err := os.Stat(name)
	return !errors.Is(err, fs.ErrNotExist)
}

// Delete deletes a key from the store
func (slf *Store) Delete(key string) error {
	pathKey := slf.PathTransform(key)
	defer func() {
		log.Printf("deleting: %s", pathKey.FullPath())
	}()

	return os.RemoveAll(path.Join(slf.Root, pathKey.FirstPathName()))
}

// Read reads a key from the store
func (slf *Store) Read(key string) (io.Reader, error) {
	fp, err := slf.readStream(key)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, fp)
	return buf, err
}

func (slf *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := slf.PathTransform(key)
	return os.Open(path.Join(slf.Root, pathKey.FullPath()))
}

func (slf *Store) writeStream(key string, reader io.Reader) error {
	pathKey := slf.PathTransform(key)
	if err := os.MkdirAll(path.Join(slf.Root, pathKey.PathName), os.ModePerm); err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, reader); err != nil {
		return err
	}

	//filenameBytes := md5.Sum(buf.Bytes())
	//filename := hex.EncodeToString(filenameBytes[:])
	//pathAndFilename := pathKey.PathName + "/" + filename
	//pathAndFilename := pathKey.FullPath()

	fp, err := os.OpenFile(path.Join(slf.Root, pathKey.FullPath()), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer fp.Close()
	n, err := io.Copy(fp, buf)
	if err != nil {
		return err
	}
	log.Printf("written (%d) bytes to: %s", n, pathKey.FullPath())
	return nil
}
