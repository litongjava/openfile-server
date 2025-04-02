package myutils

import (
  "crypto/md5"
  "encoding/hex"
  "io"
  "mime/multipart"
  "os"
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

// CalculateFileMD5FromOSFile 计算传入的 *os.File 文件的 MD5 值
func CalculateFileMD5FromOSFile(file *os.File) (string, error) {
  // 重置文件指针到文件起始位置
  if _, err := file.Seek(0, 0); err != nil {
    return "", err
  }

  hash := md5.New()
  if _, err := io.Copy(hash, file); err != nil {
    return "", err
  }

  // 可选：如果后续需要继续读取文件，可以再次重置文件指针
  if _, err := file.Seek(0, 0); err != nil {
    return "", err
  }

  return hex.EncodeToString(hash.Sum(nil)), nil
}
