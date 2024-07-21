package myutils

import (
  "github.com/cloudwego/hertz/pkg/common/hlog"
  "github.com/google/uuid"
  "io"
  "mime/multipart"
  "os"
  "path/filepath"
)

// saveFile saves the uploaded file to the given path and returns the saved file path.
func SaveFile(file *multipart.FileHeader, baseDir, fold, suffix string) (string, error) {
  uploadDir := filepath.Join(baseDir, fold)
  err := os.MkdirAll(uploadDir, os.ModePerm)
  if err != nil {
    return "", err
  }

  fileName := uuid.New().String() + suffix
  filePath := filepath.Join(uploadDir, fileName)
  err = saveUploadedFile(file, filePath)
  if err != nil {
    return "", err
  }

  fullFilePath := baseDir + fold + "/" + fileName
  hlog.Info("full file path:", fullFilePath)
  return fullFilePath, nil
}

// saveUploadedFile saves the uploaded file to the specified path.
func saveUploadedFile(file *multipart.FileHeader, destination string) error {
  src, err := file.Open()
  if err != nil {
    return err
  }
  defer src.Close()

  out, err := os.Create(destination)
  if err != nil {
    return err
  }
  defer out.Close()

  _, err = io.Copy(out, src)
  return err
}
