package myutils

import (
  "crypto/md5"
  "encoding/hex"
  "github.com/cloudwego/hertz/pkg/common/hlog"
  "github.com/google/uuid"
  "io"
  "mime/multipart"
  "os"
  "path/filepath"
)

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

  fullFilePath := baseDir + fold + "/" + fileName
  hlog.Info("full file path:", fullFilePath)
  return fullFilePath, nil
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
