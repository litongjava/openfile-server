package handler

import (
  "context"
  "github.com/litongjava/openfile-server/can"
  "net/http"
  "os"
  "path/filepath"
  "strings"
  "time"

  "github.com/cloudwego/hertz/pkg/app"
  "github.com/cloudwego/hertz/pkg/common/hlog"
  "github.com/cloudwego/hertz/pkg/common/utils"
  "github.com/litongjava/openfile-server/myutils"
)

// UploadZip 支持上传 zip 压缩包，并返回解压后所有文件的地址列表。
// 业务逻辑与之前的 Upload 方法相同，返回的 JSON 格式为 {code, message, urls}
func UploadZip(ctx context.Context, reqCtx *app.RequestContext) {
  baseDir := can.DEFAULT_FILE_PATH + "/zip/" // 基础存储目录，可根据实际需要调整

  // 1. 获取上传的 zip 文件
  fileHeader, err := reqCtx.FormFile("file")
  if err != nil {
    reqCtx.JSON(http.StatusBadRequest, utils.H{
      "code":    0,
      "message": "Failed to read file: " + err.Error(),
    })
    return
  }

  // 2. 检查文件后缀是否为 .zip
  suffix := strings.ToLower(filepath.Ext(fileHeader.Filename))
  if suffix != ".zip" {
    reqCtx.JSON(http.StatusBadRequest, utils.H{
      "code":    0,
      "message": "only support zip",
    })
    return
  }

  // 3. 获取 category 参数，若不存在则使用默认值 "default"
  category, hasCategory := reqCtx.GetPostForm("category")
  if !hasCategory || category == "" {
    category = "default"
  }
  // 以 category 与当前日期构建文件夹，例如 "default/20250401"
  fold := category + "/" + time.Now().Format("20060102")

  // 4. 生成保存 zip 文件的路径
  zipFilePath, err := myutils.GenerateFilePath(baseDir, fold, suffix)
  if err != nil {
    reqCtx.JSON(200, utils.H{
      "code":    0,
      "message": "Failed to gen file path: " + err.Error(),
    })
    return
  }

  // 5. 保存上传的 zip 文件
  if err := reqCtx.SaveUploadedFile(fileHeader, zipFilePath); err != nil {
    reqCtx.JSON(200, utils.H{
      "code":    0,
      "message": "Failed to save zip file: " + err.Error(),
    })
    return
  }

  // 6. 解压 zip 文件到目标目录（例如：将 zip 文件 "xxx.zip" 解压到 "xxx" 文件夹中）
  extractedFolder := strings.TrimSuffix(zipFilePath, suffix)
  if err := myutils.Unzip(zipFilePath, extractedFolder); err != nil {
    reqCtx.JSON(http.StatusOK, utils.H{
      "code":    0,
      "message": "failed to to unzip file:" + err.Error(),
    })
    return
  }

  // 7. 遍历解压后的所有文件，执行和 Upload 方法相同的业务逻辑
  var urls []string
  err = filepath.Walk(extractedFolder, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }
    // 忽略目录，仅处理文件
    if info.IsDir() {
      return nil
    }

    // 7.1 打开文件以计算 MD5 值
    f, err := os.Open(path)
    if err != nil {
      return err
    }
    md5Sum, err := myutils.CalculateFileMD5FromOSFile(f)
    f.Close()
    if err != nil {
      return err
    }

    // 7.2 检查数据库中是否存在该文件记录
    existingURL, err := GetFilepathFromDb(md5Sum)
    if err != nil {
      return err
    }
    var finalURL string
    if existingURL != "" {
      // 文件记录存在，检查磁盘上是否真实存在该文件
      if _, err := os.Stat(existingURL); os.IsNotExist(err) {
        // 若磁盘上不存在，则异步从当前解压文件重新保存到该路径
        go func(srcPath, destPath string) {
          srcFile, err := os.Open(srcPath)
          if err != nil {
            hlog.Error("failed to openfile:", err)
            return
          }
          defer srcFile.Close()
          if err := myutils.SaveFileFromOSFile(srcFile, destPath); err != nil {
            hlog.Error("异步保存文件失败:", err)
          }
        }(path, existingURL)
      }
      finalURL = existingURL
    } else {
      // 文件不存在：生成新的文件路径
      fileSuffix := strings.ToLower(filepath.Ext(info.Name()))
      newFilePath, err := myutils.GenerateFilePath(baseDir, fold, fileSuffix)
      if err != nil {
        return err
      }
      // 将文件信息写入数据库
      if err := SaveFileInfoToDB(md5Sum, newFilePath); err != nil {
        hlog.Error("写入数据库失败:", err)
        // 即使数据库写入失败，依然继续处理其他文件
      }
      // 异步保存文件到新的位置
      go func(srcPath, destPath string) {
        srcFile, err := os.Open(srcPath)
        if err != nil {
          hlog.Error("打开文件失败:", err)
          return
        }
        defer srcFile.Close()
        if err := myutils.SaveFileFromOSFile(srcFile, destPath); err != nil {
          hlog.Error("异步保存文件失败:", err)
        }
      }(path, newFilePath)
      finalURL = newFilePath
    }

    // 7.3 添加url
    finalURL = filepath.ToSlash(finalURL)
    urls = append(urls, strings.TrimPrefix(finalURL, baseDir))
    return nil
  })
  if err != nil {
    reqCtx.JSON(http.StatusInternalServerError, utils.H{
      "code":    0,
      "message": "Faile to replace zip file: " + err.Error(),
    })
    return
  }

  urlPrefix := myutils.GetFullHostURL(reqCtx.URI())
  // 8. 返回上传成功信息及所有解压后文件的 URL 列表
  reqCtx.JSON(http.StatusOK, utils.H{
    "code":   200,
    "urls":   urls,
    "server": urlPrefix,
  })
}
