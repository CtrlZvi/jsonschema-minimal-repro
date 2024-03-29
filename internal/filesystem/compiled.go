// Code generated by fileb0x at "2019-08-22 15:56:49.2371968 -0700 PDT m=+0.012998901" from config file "b0x.yaml" DO NOT EDIT.
// modification hash(dea31b30019963c58b1e466dea82b04d.e783fc83d1ef42c16a3dbf2e923a4267)

package filesystem

import (
	"bytes"

	"context"
	"io"
	"net/http"
	"os"
	"path"

	"golang.org/x/net/webdav"
)

var (
	// CTX is a context for webdav vfs
	CTX = context.Background()

	// FS is a virtual memory file system
	FS = webdav.NewMemFS()

	// Handler is used to server files through a http handler
	Handler *webdav.Handler

	// HTTP is the http file system
	HTTP http.FileSystem = new(HTTPFS)
)

// HTTPFS implements http.FileSystem
type HTTPFS struct {
	// Prefix allows to limit the path of all requests. F.e. a prefix "css" would allow only calls to /css/*
	Prefix string
}

// FileSchemasSchemaJSON is "schemas/schema.json"
var FileSchemasSchemaJSON = []byte("\x7b\x0d\x0a\x20\x20\x20\x20\x22\x24\x73\x63\x68\x65\x6d\x61\x22\x3a\x20\x22\x68\x74\x74\x70\x3a\x2f\x2f\x6a\x73\x6f\x6e\x2d\x73\x63\x68\x65\x6d\x61\x2e\x6f\x72\x67\x2f\x64\x72\x61\x66\x74\x2d\x30\x37\x2f\x73\x63\x68\x65\x6d\x61\x23\x22\x2c\x0d\x0a\x20\x20\x20\x20\x22\x74\x79\x70\x65\x22\x3a\x20\x22\x6f\x62\x6a\x65\x63\x74\x22\x2c\x0d\x0a\x20\x20\x20\x20\x22\x74\x69\x74\x6c\x65\x22\x3a\x20\x22\x4d\x69\x6e\x69\x6d\x75\x6d\x20\x52\x65\x70\x72\x6f\x20\x43\x61\x73\x65\x22\x2c\x0d\x0a\x20\x20\x20\x20\x22\x64\x65\x73\x63\x72\x69\x70\x74\x69\x6f\x6e\x22\x3a\x20\x22\x41\x6e\x20\x65\x78\x61\x6d\x70\x6c\x65\x20\x73\x63\x68\x65\x6d\x61\x20\x66\x6f\x72\x20\x75\x73\x65\x20\x69\x6e\x20\x6d\x69\x6e\x69\x6d\x75\x6d\x20\x72\x65\x70\x72\x6f\x20\x63\x61\x73\x65\x73\x22\x2c\x0d\x0a\x7d")

func init() {
	err := CTX.Err()
	if err != nil {
		panic(err)
	}

	err = FS.Mkdir(CTX, "schemas/", 0777)
	if err != nil && err != os.ErrExist {
		panic(err)
	}

	var f webdav.File

	f, err = FS.OpenFile(CTX, "schemas/schema.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}

	_, err = f.Write(FileSchemasSchemaJSON)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

	Handler = &webdav.Handler{
		FileSystem: FS,
		LockSystem: webdav.NewMemLS(),
	}

}

// Open a file
func (hfs *HTTPFS) Open(path string) (http.File, error) {
	path = hfs.Prefix + path

	f, err := FS.OpenFile(CTX, path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// ReadFile is adapTed from ioutil
func ReadFile(path string) ([]byte, error) {
	f, err := FS.OpenFile(CTX, path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(make([]byte, 0, bytes.MinRead))

	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(f)
	return buf.Bytes(), err
}

// WriteFile is adapTed from ioutil
func WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := FS.OpenFile(CTX, filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

// WalkDirs looks for files in the given dir and returns a list of files in it
// usage for all files in the b0x: WalkDirs("", false)
func WalkDirs(name string, includeDirsInList bool, files ...string) ([]string, error) {
	f, err := FS.OpenFile(CTX, name, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	fileInfos, err := f.Readdir(0)
	if err != nil {
		return nil, err
	}

	err = f.Close()
	if err != nil {
		return nil, err
	}

	for _, info := range fileInfos {
		filename := path.Join(name, info.Name())

		if includeDirsInList || !info.IsDir() {
			files = append(files, filename)
		}

		if info.IsDir() {
			files, err = WalkDirs(filename, includeDirsInList, files...)
			if err != nil {
				return nil, err
			}
		}
	}

	return files, nil
}
