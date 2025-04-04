package handler

import (
  "context"
  "github.com/litongjava/openfile-server/can"
  "net/http"
  "path/filepath"
  "strings"
  "time"

  "github.com/cloudwego/hertz/pkg/app"
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
  zipFilePath, err := myutils.GenerateFilePathWithName(baseDir, fold, fileHeader.Filename)
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
  urls, err := myutils.Unzip(zipFilePath, extractedFolder)
  if err != nil {
    reqCtx.JSON(http.StatusOK, utils.H{
      "code":    0,
      "message": "failed to to unzip file:" + err.Error(),
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
