package handler

import (
  "fmt"
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
  filePath, err := GetExistingFileURL(md5Sum)
  if err == nil && filePath != "" {
    _, err := os.Stat(filePath)
    if !os.IsNotExist(err) {
      hlog.Info("file exists")
    }
  } else {
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
    go func() {
      err := myutils.SaveFile(file, filePath)
      if err != nil {
        hlog.Error("Failed to save file:", err)
      }
    }()
  }

  // 检查是否为视频文件
  isVideo := false
  videoExtensions := map[string]bool{
    ".mp4": true, ".avi": true, ".mov": true, ".mkv": true, ".flv": true,
  }
  if videoExtensions[suffix] {
    isVideo = true
  }

  var frames []string
  if isVideo {
    // 获取视频时长
    duration, err := myutils.GetVideoDuration(filePath)
    if err != nil {
      hlog.Error("Failed to get video duration:", err)
    } else {
      var frameCount int
      if duration >= 10 {
        frameCount = 10
      } else {
        frameCount = int(duration)
        if frameCount < 1 {
          frameCount = 1
        }
      }

      // 提取关键帧
      frameDir := filepath.Join("file", "frames", fold)
      framePaths, err := myutils.ExtractKeyFrames(filePath, frameDir, frameCount)
      if err != nil {
        hlog.Error("Failed to extract key frames:", err)
      } else {
        for _, framePath := range framePaths {
          // 生成雪花ID作为文件名
          snowflakeID := myutils.GenerateSnowflakeID()
          newFrameFilename := fmt.Sprintf("%s%s", snowflakeID, filepath.Ext(framePath))
          newFramePath := filepath.Join(filepath.Dir(framePath), newFrameFilename)

          // 重命名帧文件，确保使用正斜杠
          err := os.Rename(framePath, newFramePath)
          if err != nil {
            hlog.Error("Failed to rename frame file:", err)
            continue
          }
          var relativeFramePath = filepath.ToSlash(newFramePath)

          // 添加到 frames 列表
          frames = append(frames, relativeFramePath)
        }
      }
    }
  }

  // 构建响应
  response := UploadVideoResponse{
    Code:   200,
    URL:    urlPrefix,
    Data:   filePath,
    MD5:    md5Sum,
    Frames: frames,
  }

  reqCtx.JSON(http.StatusOK, response)
}
