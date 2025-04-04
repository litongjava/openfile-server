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

  // 获取存储目录
  category, hasCategory := reqCtx.GetPostForm("category")
  if !hasCategory || category == "" {
    category = "default"
  }
  fold := category + "/" + time.Now().Format("20060102")

  // 生成保存上传分片文件的路径
  shardFilePath, err := myutils.GenerateFilePathWithName(baseDir, fold, fileName+"."+strconv.Itoa(partIndex))
  if err != nil {
    reqCtx.JSON(200, utils.H{
      "code":    0,
      "message": "Failed to gen file path: " + err.Error(),
    })
    return
  }

  // 保存上传的分片文件
  if err := reqCtx.SaveUploadedFile(fileHeader, shardFilePath); err != nil {
    reqCtx.JSON(200, utils.H{
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

// 合并所有分片并解压
func mergeShardsAndUnzip(baseDir, fold, fileName string, totalParts int) ([]string, error) {
  var fileData []byte
  for i := 0; i < totalParts; i++ {
    shardPath := baseDir + fold + "/" + fileName + "." + strconv.Itoa(i)
    shardData, err := os.ReadFile(shardPath)
    if err != nil {
      hlog.Errorf("failed to read shard %d: %v", i, err)
      return nil, fmt.Errorf("failed to read shard %d: %v", i, err)
    }
    fileData = append(fileData, shardData...)
  }

  // 保存合并后的文件
  completeFilePath := baseDir + fold + "/" + fileName
  if err := os.WriteFile(completeFilePath, fileData, 0644); err != nil {
    return nil, fmt.Errorf("failed to save merged file: %v", err)
  }

  // 解压文件
  extractedFolder := strings.TrimSuffix(completeFilePath, ".zip")
  return myutils.Unzip(completeFilePath, extractedFolder)
}
