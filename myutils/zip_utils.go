package myutils

import (
  "archive/zip"
  "fmt"
  "io"
  "os"
  "path/filepath"
  "strings"
)

// Unzip 解压 zipFilePath 指定的压缩包到 extractedFolder 目录下，并返回所有解压后的文件列表，格式如 file/service/20250403/1/xx.png
func Unzip(zipFilePath, extractedFolder string) ([]string, error) {
  // 打开 zip 文件
  r, err := zip.OpenReader(zipFilePath)
  if err != nil {
    return nil, fmt.Errorf("打开压缩包失败: %w", err)
  }
  defer r.Close()

  // 确保目标目录存在
  err = os.MkdirAll(extractedFolder, os.ModePerm)
  if err != nil {
    return nil, fmt.Errorf("创建解压目录失败: %w", err)
  }

  var extractedFiles []string

  // 遍历 zip 内所有的文件/目录
  for _, f := range r.File {
    // 生成每个文件/目录的目标路径
    fpath := filepath.Join(extractedFolder, f.Name)

    // 防止 ZipSlip 漏洞，确保解压路径不会跳出目标目录
    if !strings.HasPrefix(fpath, filepath.Clean(extractedFolder)+string(os.PathSeparator)) {
      return nil, fmt.Errorf("非法文件路径: %s", fpath)
    }

    if f.FileInfo().IsDir() {
      // 如果是目录，则创建目录
      if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
        return nil, fmt.Errorf("创建目录失败: %w", err)
      }
      continue
    }

    // 确保文件所在目录存在
    if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
      return nil, fmt.Errorf("创建文件目录失败: %w", err)
    }

    // 创建目标文件
    outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
    if err != nil {
      return nil, fmt.Errorf("创建目标文件失败: %w", err)
    }

    // 打开 zip 内的文件
    rc, err := f.Open()
    if err != nil {
      outFile.Close()
      return nil, fmt.Errorf("打开压缩包内文件失败: %w", err)
    }

    // 将文件内容复制到目标文件
    _, err = io.Copy(outFile, rc)
    // 关闭文件
    outFile.Close()
    rc.Close()

    if err != nil {
      return nil, fmt.Errorf("写入文件失败: %w", err)
    }

    // 转换路径为标准的正斜杠格式，并添加到列表中
    normalizedPath := strings.ReplaceAll(fpath, string(os.PathSeparator), "/")
    extractedFiles = append(extractedFiles, normalizedPath)
  }

  return extractedFiles, nil
}
