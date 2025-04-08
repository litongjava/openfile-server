package handler

import (
  "context"
  "net/http"
  "path/filepath"
  "strings"
  "time"

  "github.com/cloudwego/hertz/pkg/app"
  "github.com/cloudwego/hertz/pkg/common/hlog"
  "github.com/cloudwego/hertz/pkg/common/utils"
  "github.com/litongjava/openfile-server/can"
  "github.com/litongjava/openfile-server/myutils"
)

// UploadZip 支持上传 zip 压缩包，并返回解压后所有文件的地址列表。
// 新增逻辑：压缩包解压完成后，判断是否有视频文件，如果有则将其转为 HLS 流（转换操作采用异步方式）。
func UploadZip(ctx context.Context, reqCtx *app.RequestContext) {
  // 设置基础存储目录，此处为 zip 类型文件的存储目录，根据实际需要可调整
  baseDir := can.DEFAULT_FILE_PATH + "/zip/"

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
  // 构造保存路径中的文件夹名称，例如 "default/20250401"
  fold := category + "/" + time.Now().Format("20060102")

  // 4. 生成保存 zip 文件的路径，使用自定义的 GenerateFilePathWithName 方法
  zipFilePath, err := myutils.GenerateFilePathWithName(baseDir, fold, fileHeader.Filename)
  if err != nil {
    reqCtx.JSON(http.StatusOK, utils.H{
      "code":    0,
      "message": "Failed to gen file path: " + err.Error(),
    })
    return
  }

  // 5. 保存上传的 zip 文件
  if err := reqCtx.SaveUploadedFile(fileHeader, zipFilePath); err != nil {
    reqCtx.JSON(http.StatusOK, utils.H{
      "code":    0,
      "message": "Failed to save zip file: " + err.Error(),
    })
    return
  }

  // 6. 解压 zip 文件到目标目录（例如将 "xxx.zip" 解压到 "xxx" 文件夹中）
  // extractedFolder 为 zip 文件去掉后缀后的文件夹名称
  extractedFolder := strings.TrimSuffix(zipFilePath, suffix)
  urls, err := myutils.Unzip(zipFilePath, extractedFolder)
  if err != nil {
    reqCtx.JSON(http.StatusOK, utils.H{
      "code":    0,
      "message": "failed to unzip file: " + err.Error(),
    })
    return
  }

  // 7. 新增逻辑：
  //    遍历解压后的文件列表，判断是否有视频文件，如果有则调用 ConvertVideoToHLS 转换为 HLS 流。
  videoExtensions := map[string]bool{
    ".mp4": true,
    ".avi": true,
    ".mov": true,
    ".mkv": true,
    ".flv": true,
  }
  for _, filePath := range urls {
    ext := strings.ToLower(filepath.Ext(filePath))
    if videoExtensions[ext] {
      // 异步调用转换，不阻塞返回数据
      go func(fp, ext string) {
        _, err := myutils.ConvertVideoToHLS(fp, baseDir, ext)
        if err != nil {
          hlog.Error("HLS conversion failed for", fp, ":", err)
        } else {
          hlog.Info("HLS conversion succeeded for", fp)
        }
      }(filePath, ext)
    }
  }

  // 8. 构造返回数据，返回解压后所有文件的 URL 列表及服务器前缀信息，返回数据结构保持不变
  urlPrefix := myutils.GetFullHostURL(reqCtx.URI())
  reqCtx.JSON(http.StatusOK, utils.H{
    "code":   200,
    "urls":   urls,
    "server": urlPrefix,
  })
}
