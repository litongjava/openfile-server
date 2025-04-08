package handler

import (
  "context"
  "fmt"
  "github.com/cloudwego/hertz/pkg/app"
  "github.com/cloudwego/hertz/pkg/common/hlog"
  "github.com/cloudwego/hertz/pkg/common/utils"
  "github.com/litongjava/openfile-server/can"
  "github.com/litongjava/openfile-server/myutils"
  "net/http"
  "os"
  "path/filepath"
  "strconv"
  "strings"
  "time"
)

// UploadZipShard 支持分片上传 zip 压缩包
func UploadZipShard(ctx context.Context, reqCtx *app.RequestContext) {
  baseDir := can.DEFAULT_FILE_PATH + "/zip/" // 基础存储目录

  // 获取上传的分片文件
  fileHeader, err := reqCtx.FormFile("file")
  if err != nil {
    reqCtx.JSON(http.StatusBadRequest, utils.H{
      "code":    0,
      "message": "Failed to read file: " + err.Error(),
    })
    return
  }

  // 获取分片的编号与总分片数
  partIndexStr := reqCtx.DefaultPostForm("partIndex", "0")
  totalPartsStr := reqCtx.DefaultPostForm("totalParts", "0")

  // 将字符串转换为整数
  partIndex, err := strconv.Atoi(partIndexStr)
  if err != nil {
    reqCtx.JSON(http.StatusBadRequest, utils.H{
      "code":    0,
      "message": "Invalid partIndex value: " + partIndexStr,
    })
    return
  }

  totalParts, err := strconv.Atoi(totalPartsStr)
  if err != nil {
    reqCtx.JSON(http.StatusBadRequest, utils.H{
      "code":    0,
      "message": "Invalid totalParts value: " + totalPartsStr,
    })
    return
  }

  // 获取文件名称和文件后缀
  fileName := reqCtx.DefaultPostForm("fileName", "")
  suffix := filepath.Ext(fileName)

  if suffix != ".zip" {
    reqCtx.JSON(http.StatusBadRequest, utils.H{
      "code":    0,
      "message": "only support zip",
    })
    return
  }

  // 获取存储目录，如果没有指定则默认 "default"
  category, hasCategory := reqCtx.GetPostForm("category")
  if !hasCategory || category == "" {
    category = "default"
  }
  // 以 category 与当前日期构建文件夹，例如 "default/20250401"
  fold := category + "/" + time.Now().Format("20060102")

  // 生成保存上传分片文件的路径
  shardFilePath, err := myutils.GenerateFilePathWithName(baseDir, fold, fileName+"."+strconv.Itoa(partIndex))
  if err != nil {
    reqCtx.JSON(http.StatusOK, utils.H{
      "code":    0,
      "message": "Failed to gen file path: " + err.Error(),
    })
    return
  }

  // 保存上传的分片文件
  if err := reqCtx.SaveUploadedFile(fileHeader, shardFilePath); err != nil {
    reqCtx.JSON(http.StatusOK, utils.H{
      "code":    0,
      "message": "Failed to save shard file: " + err.Error(),
    })
    return
  }

  // 判断是否是最后一个分片
  if partIndex == totalParts-1 {
    // 合并所有分片并解压
    urls, err := mergeShardsAndUnzip(baseDir, fold, fileName, totalParts)
    if err != nil {
      reqCtx.JSON(http.StatusOK, utils.H{
        "code":    0,
        "message": "Failed to merge and unzip: " + err.Error(),
      })
      return
    }

    // 新增逻辑：遍历解压后的文件列表，如果包含视频文件，异步转换成 HLS 流格式
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
        // 异步调用转换，调用 ConvertVideoToHLS 将视频文件切片成 HLS 格式
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

    urlPrefix := myutils.GetFullHostURL(reqCtx.URI())
    // 返回上传成功信息及所有解压后文件的 URL 列表
    reqCtx.JSON(http.StatusOK, utils.H{
      "code":   200,
      "urls":   urls,
      "server": urlPrefix,
    })
  } else {
    reqCtx.JSON(http.StatusOK, utils.H{
      "code":    200,
      "message": "Shard uploaded successfully",
    })
  }
}

// mergeShardsAndUnzip 合并所有分片并解压，同时删除服务器中的分片文件
func mergeShardsAndUnzip(baseDir, fold, fileName string, totalParts int) ([]string, error) {
  var fileData []byte
  // 用来记录所有分片文件的路径，便于后续删除
  var shardPaths []string

  for i := 0; i < totalParts; i++ {
    shardPath := baseDir + fold + "/" + fileName + "." + strconv.Itoa(i)
    shardPaths = append(shardPaths, shardPath)
    shardData, err := os.ReadFile(shardPath)
    if err != nil {
      hlog.Errorf("failed to read shard %d: %v", i, err)
      return nil, fmt.Errorf("failed to read shard %d: %v", i, err)
    }
    fileData = append(fileData, shardData...)
  }

  // 保存合并后的完整文件
  completeFilePath := baseDir + fold + "/" + fileName
  if err := os.WriteFile(completeFilePath, fileData, 0644); err != nil {
    return nil, fmt.Errorf("failed to save merged file: %v", err)
  }

  // 删除所有分片文件
  for _, shardPath := range shardPaths {
    if err := os.Remove(shardPath); err != nil {
      hlog.Errorf("failed to delete shard file %s: %v", shardPath, err)
    }
  }

  // 解压文件，将 zip 文件解压到去掉后缀的文件夹下
  extractedFolder := strings.TrimSuffix(completeFilePath, ".zip")
  return myutils.Unzip(completeFilePath, extractedFolder)
}
