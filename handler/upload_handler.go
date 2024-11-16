package handler

import (
  "context"
  "fmt"
  "github.com/cloudwego/hertz/pkg/app"
  "github.com/cloudwego/hertz/pkg/common/hlog"
  "github.com/cloudwego/hertz/pkg/common/utils"
  "github.com/cloudwego/hertz/pkg/protocol/consts"
  "github.com/litongjava/openfile-server/myutils"
  "log"
  "net/http"
  "os"
  "path/filepath"
  "strings"
  "time"
)

func UploadHandler(ctx context.Context, reqCtx *app.RequestContext) {
  username := reqCtx.Param("username")             // 获取url中的username
  repositoryName := reqCtx.Param("repositoryName") // 获取url中的repositoryName
  subFolder := reqCtx.Param("subFolder")           // 获取url中的filePath

  fileHeader, err := reqCtx.FormFile("file") // 获取上传的文件
  if err != nil {
    reqCtx.JSON(consts.StatusBadRequest, utils.H{
      "error": err.Error(),
    })
    return
  }

  var uploadDir = filepath.Join("s", username, repositoryName, subFolder)
  err = os.MkdirAll(uploadDir, os.ModePerm)
  if err != nil {
    reqCtx.JSON(consts.StatusInternalServerError, utils.H{
      "error": err.Error(),
    })
  }

  filename := fileHeader.Filename
  path := filepath.Join(uploadDir, filename) // 构建保存的文件路径
  log.Println("path:", path)

  if err := reqCtx.SaveUploadedFile(fileHeader, path); err != nil { // 保存文件
    reqCtx.JSON(consts.StatusInternalServerError, utils.H{
      "error": err.Error(),
    })
    return
  }

  reqCtx.JSON(consts.StatusOK, utils.H{
    "message": fmt.Sprintf("'%s' uploaded!", filename),
  })
}

// Upload handles file uploads and returns the file path.
func Upload(reqCtx *app.RequestContext, baseDir string) {
  file, err := reqCtx.FormFile("file")
  if err != nil {
    reqCtx.JSON(http.StatusBadRequest, utils.H{
      "code": 0,
      "data": "Failed to read file",
    })
    return
  }

  category, hasCategory := reqCtx.GetPostForm("category")
  fold := category
  if !hasCategory {
    fold = time.Now().Format("20060102")
  }

  suffix := strings.ToLower(filepath.Ext(file.Filename))
  md5Sum, err := myutils.CalculateFileMD5(file)
  if err != nil {
    reqCtx.JSON(http.StatusInternalServerError, utils.H{
      "code": 0,
      "data": "Failed to calculate file MD5",
    })
    return
  }

  // Check if file already exists in DB
  existingURL, err := GetExistingFileURL(md5Sum)
  if err == nil && existingURL != "" {
    _, err := os.Stat(existingURL)
    if os.IsNotExist(err) {
      go func() {
        err := myutils.SaveFile(file, existingURL)
        if err != nil {
          hlog.Error("Failed to save file:", err)
        } else {
          hlog.Info("old file save success")
        }
      }()
    } else {
      hlog.Info("file exists")
    }
    reqCtx.JSON(http.StatusOK, utils.H{
      "code": 200,
      "url":  myutils.GetFullHostURL(reqCtx.URI()),
      "data": existingURL,
      "md5":  md5Sum,
    })
    return
  }

  filePath, err := myutils.GenerateFilePath(baseDir, fold, suffix)
  if err != nil {
    reqCtx.JSON(http.StatusInternalServerError, utils.H{
      "code": 0,
      "data": "Failed to generate file path",
    })
    return
  }

  url := myutils.GetFullHostURL(reqCtx.URI())
  err = SaveFileInfoToDB(md5Sum, filePath)
  if err != nil {
    reqCtx.JSON(http.StatusOK, utils.H{
      "code":  -1,
      "url":   url,
      "data":  filePath,
      "md5":   md5Sum,
      "error": err.Error(),
    })
    return
  }

  go func() {
    err := myutils.SaveFile(file, filePath)
    if err != nil {
      hlog.Error("Failed to save file:", err)
    }
  }()

  reqCtx.JSON(http.StatusOK, utils.H{
    "code": 200,
    "url":  url,
    "data": filePath,
    "md5":  md5Sum,
  })
}
