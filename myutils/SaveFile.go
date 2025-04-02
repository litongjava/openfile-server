package myutils

import (
  "github.com/cloudwego/hertz/pkg/common/hlog"
  "github.com/google/uuid"
  "io"
  "mime/multipart"
  "os"
  "path/filepath"
)

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

// SaveFileFromOSFile 将传入的 *os.File 文件内容保存到指定的 filePath 路径
func SaveFileFromOSFile(file *os.File, filePath string) error {
  // 确保目标目录存在
  dir := filepath.Dir(filePath)
  if _, err := os.Stat(dir); os.IsNotExist(err) {
    // 创建目录
    if err := os.MkdirAll(dir, os.ModePerm); err != nil {
      return err
    }
  }

  // 可选：重置文件指针到起始位置，确保从头开始复制
  if _, err := file.Seek(0, 0); err != nil {
    return err
  }

  // 创建目标文件
  dst, err := os.Create(filePath)
  if err != nil {
    return err
  }
  defer dst.Close()

  // 复制文件内容到目标文件
  _, err = io.Copy(dst, file)
  return err
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
