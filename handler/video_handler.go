package handler

import (
  "context"
  "github.com/cloudwego/hertz/pkg/app"
  "github.com/cloudwego/hertz/pkg/common/utils"
  "net/http"
  "os"
  "strings"
  "time"
)

func VideoFrames(ctx context.Context, reqCtx *app.RequestContext) {
  filePath := reqCtx.Query("uri")

  var frameArray []string
  if filePath != "" {
    _, err := os.Stat(filePath)
    if os.IsNotExist(err) {
      reqCtx.String(400, "file is not exists")
      return
    }
    err, framesString := GetVideoFramesFromDb(filePath)
    if err != nil {
      reqCtx.JSON(http.StatusInternalServerError, utils.H{
        "code": 0,
        "data": err.Error(),
      })
      return
    }
    frameArray = strings.Split(framesString, ",")
  } else {
    md5Sum, err := GetMd5ByFiepath(filePath)
    if err != nil {
      reqCtx.JSON(400, utils.H{
        "code": 0,
        "data": err.Error(),
      })
      return
    }
    var fold = time.Now().Format("20060102")
    frameArray = ExtraFrames(filePath, fold)
    result := strings.Join(frameArray, ",")
    SaveVideoFramesToDB(md5Sum, filePath, result)
  }

  // 构建响应
  response := UploadVideoResponse{
    Code:   200,
    Data:   filePath,
    Frames: frameArray,
  }
  reqCtx.JSON(http.StatusOK, response)
}
