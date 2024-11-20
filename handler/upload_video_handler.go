package handler

import (
  "github.com/cloudwego/hertz/pkg/app"
  "github.com/cloudwego/hertz/pkg/common/hlog"
  "github.com/cloudwego/hertz/pkg/common/utils"
  "github.com/litongjava/openfile-server/myutils"
  "net/http"
  "os"
  "path/filepath"
  "strings"
  "time"
)

// UploadVideoResponse 定义视频上传响应结构
type UploadVideoResponse struct {
  Code   int      `json:"code"`
  Data   string   `json:"data"`
  URL    string   `json:"url"`
  MD5    string   `json:"md5"`
  Frames []string `json:"frames,omitempty"`
}

// UploadVideo 处理视频文件的上传
func UploadVideo(reqCtx *app.RequestContext, baseDir string) {
  // 获取上传的文件
  file, err := reqCtx.FormFile("file")
  if err != nil {
    reqCtx.JSON(http.StatusBadRequest, utils.H{
      "code": 0,
      "data": "Failed to read file",
    })
    return
  }

  // 获取分类（如果有）
  category, hasCategory := reqCtx.GetPostForm("category")
  fold := category
  if !hasCategory {
    fold = time.Now().Format("20060102")
  }

  // 获取文件后缀名并转换为小写
  suffix := strings.ToLower(filepath.Ext(file.Filename))

  // 计算文件 MD5
  md5Sum, err := myutils.CalculateFileMD5(file)
  if err != nil {
    reqCtx.JSON(http.StatusInternalServerError, utils.H{
      "code": 0,
      "data": "Failed to calculate file MD5",
    })
    return
  }

  // 获取服务器的完整 URL 前缀
  urlPrefix := myutils.GetFullHostURL(reqCtx.URI())

  // Check if file already exists in DB
  filePath, err := GetFilepathFromDb(md5Sum)
  var frameArray []string
  if err == nil && filePath != "" {
    _, err := os.Stat(filePath)
    if !os.IsNotExist(err) {
      hlog.Info("file exists")
    } else {
      // 保存主文件
      err := myutils.SaveFile(file, filePath)
      if err != nil {
        hlog.Error("Failed to save file:", err)
      }
    }
    err, framesString := GetVideoFramesFromDb(filePath)
    if err != nil {
      reqCtx.JSON(http.StatusInternalServerError, utils.H{
        "code": 0,
        "data": err.Error(),
      })
      return
    }
    if framesString != "" {
      frameArray = strings.Split(framesString, ",")
    } else {
      frameArray = ExtraFrames(filePath, fold)
      result := strings.Join(frameArray, ",")
      SaveVideoFramesToDB(md5Sum, filePath, result)
    }

    // 构建响应
    response := UploadVideoResponse{
      Code:   200,
      URL:    urlPrefix,
      Data:   filePath,
      MD5:    md5Sum,
      Frames: frameArray,
    }
    reqCtx.JSON(http.StatusOK, response)
    return
  }
  // 生成文件保存路径
  filePath, err = myutils.GenerateFilePath(baseDir, fold, suffix)
  if err != nil {
    reqCtx.JSON(http.StatusInternalServerError, utils.H{
      "code": 0,
      "data": "Failed to generate file path",
    })
    return
  }

  // 保存文件信息到数据库
  err = SaveFileInfoToDB(md5Sum, filePath)
  if err != nil {
    reqCtx.JSON(http.StatusOK, utils.H{
      "code":  -1,
      "url":   urlPrefix,
      "data":  filePath,
      "md5":   md5Sum,
      "error": err.Error(),
    })
    return
  }
  // 保存主文件
  err = myutils.SaveFile(file, filePath)
  if err != nil {
    hlog.Error("Failed to save file:", err)
  }
  // 检查是否为视频文件
  isVideo := false
  videoExtensions := map[string]bool{
    ".mp4": true, ".avi": true, ".mov": true, ".mkv": true, ".flv": true,
  }
  if videoExtensions[suffix] {
    isVideo = true
  }
  if isVideo {
    frameArray = ExtraFrames(filePath, fold)
    result := strings.Join(frameArray, ",")
    SaveVideoFramesToDB(md5Sum, filePath, result)
  }

  // 构建响应
  response := UploadVideoResponse{
    Code:   200,
    URL:    urlPrefix,
    Data:   filePath,
    MD5:    md5Sum,
    Frames: frameArray,
  }

  reqCtx.JSON(http.StatusOK, response)
}
