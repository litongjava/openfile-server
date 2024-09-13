package myutils

import (
  "fmt"
  "image"
  "image/gif"
  "image/jpeg"
  "image/png"
  "os"
  "path/filepath"
  "strings"

  "github.com/chai2010/webp"
  "github.com/nfnt/resize"
  "golang.org/x/image/bmp"
  "golang.org/x/image/tiff"
)

func GenerateThumbnail(originalPath, thumbnailPath string, width, height int) error {
  file, err := os.Open(originalPath)
  if err != nil {
    return fmt.Errorf("failed to open original file: %v", err)
  }
  defer file.Close()

  // 解码图片
  img, format, err := image.Decode(file)
  if err != nil {
    return fmt.Errorf("failed to decode image: %v", err)
  }

  // 生成缩略图
  thumb := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

  // 获取缩略图文件的目录
  dir := filepath.Dir(thumbnailPath)

  // 检查目录是否存在，如果不存在则创建
  if _, err := os.Stat(dir); os.IsNotExist(err) {
    err = os.MkdirAll(dir, os.ModePerm)
    if err != nil {
      return fmt.Errorf("failed to create directory: %v", err)
    }
  }
  // 创建缩略图文件
  out, err := os.Create(thumbnailPath)
  if err != nil {
    return fmt.Errorf("failed to create thumbnail file: %v", err)
  }
  defer out.Close()

  // 根据格式保存图片
  switch strings.ToLower(format) {
  case "jpeg", "jpg":
    err = jpeg.Encode(out, thumb, nil)
  case "png":
    err = png.Encode(out, thumb)
  case "gif":
    err = gif.Encode(out, thumb, nil)
  case "bmp":
    err = bmp.Encode(out, thumb)
  case "tiff":
    err = tiff.Encode(out, thumb, nil)
  case "webp":
    err = webp.Encode(out, thumb, nil)
  default:
    return fmt.Errorf("unsupported image format: %s", format)
  }

  if err != nil {
    return fmt.Errorf("failed to save thumbnail: %v", err)
  }

  return nil
}
