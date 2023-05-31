package cfg

import (
   "bytes"
   "errors"
   "io/fs"
   "os"
   "path"
)

type FileSystem interface {
   Exists(filepath string) (bool, error)
   Read(filepath string) ([]byte, error)
}

type defaultFs struct {
   fs   fs.FS
   root string
}

func defaultFileSystem() (*defaultFs, error) {
   if pwd, err := os.Getwd(); err != nil {
      return nil, err
   } else {
      return &defaultFs{
         root: pwd,
         fs:   os.DirFS(pwd),
      }, nil
   }
}

func (f *defaultFs) open(filepath string) (fs.File, error) {
   fullPath := path.Join(f.root, filepath)
   if file, err := f.fs.Open(fullPath); err != nil {
      if pathError, isPathError := err.(*fs.PathError); isPathError && pathError.Err == fs.ErrNotExist {
         return nil, fs.ErrNotExist
      }
      return nil, err
   } else if file != nil {
      return file, nil
   }
   return nil, fs.ErrNotExist
}

func (f *defaultFs) Exists(filepath string) (bool, error) {
   if file, err := f.open(filepath); err != nil {
      if errors.Is(err, os.ErrNotExist) {
         return false, nil
      }
      return false, err
   } else {
      _ = file.Close()
      return true, nil
   }
}

func (f *defaultFs) Read(filepath string) ([]byte, error) {
   if file, err := f.open(filepath); err != nil {
      return nil, err
   } else {
      defer func(file fs.File) {
         _ = file.Close()
      }(file)

      buf := new(bytes.Buffer)
      if _, err = buf.ReadFrom(file); err != nil {
         return nil, err
      }
      return buf.Bytes(), nil
   }
}
