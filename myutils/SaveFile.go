package myutils

import (
  "crypto/md5"
  "encoding/hex"
  "github.com/google/uuid"
  "hash"
  "io"
  "mime/multipart"
  "os"
  "path/filepath"
)

// saveUploadedFile saves the uploaded file to the specified path.
func saveUploadedFile(file *multipart.FileHeader, destination string) (hash.Hash, error) {
  src, err := file.Open()
  if err != nil {
    return nil, err
  }
  defer src.Close()

  out, err := os.Create(destination)
  if err != nil {
    return nil, err
  }
  defer out.Close()

  hash := md5.New()
  // TeeReader will write to both the out and the hash
  tee := io.TeeReader(src, out)

  // Copy the file while calculating the MD5 hash
  if _, err := io.Copy(hash, tee); err != nil {
    return nil, err
  }

  return hash, nil
}

func CalculateFileMD5(file *multipart.FileHeader) (string, error) {
  src, err := file.Open()
  if err != nil {
    return "", err
  }
  defer src.Close()

  hash := md5.New()
  if _, err := io.Copy(hash, src); err != nil {
    return "", err
  }

  return hex.EncodeToString(hash.Sum(nil)), nil
}

func GenerateFilePath(baseDir, fold, suffix string) (string, error) {
  uploadDir := filepath.Join(baseDir, fold)
  err := os.MkdirAll(uploadDir, os.ModePerm)
  if err != nil {
    return "", err
  }

  fileName := uuid.New().String() + suffix
  return filepath.Join(fold, fileName), nil
}

func SaveFile(file *multipart.FileHeader, filePath string) error {
  src, err := file.Open()
  if err != nil {
    return err
  }
  defer src.Close()

  dst, err := os.Create(filePath)
  if err != nil {
    return err
  }
  defer dst.Close()

  _, err = io.Copy(dst, src)
  return err
}
