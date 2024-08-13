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
  // 获取目录路径
  dir := filepath.Dir(filePath)

  // 检查目录是否存在
  if _, err := os.Stat(dir); os.IsNotExist(err) {
    // 创建目录
    err := os.MkdirAll(dir, os.ModePerm)
    if err != nil {
      return err
    }
  }

  // 打开源文件
  src, err := file.Open()
  if err != nil {
    return err
  }
  defer src.Close()

  // 创建目标文件
  dst, err := os.Create(filePath)
  if err != nil {
    return err
  }
  defer dst.Close()

  // 复制文件内容
  _, err = io.Copy(dst, src)
  return err
}
