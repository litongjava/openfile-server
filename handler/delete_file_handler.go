package handler

import (
  "context"
  "github.com/cloudwego/hertz/pkg/app"
  "github.com/cloudwego/hertz/pkg/common/utils"
  "net/http"
  "os"
  "strings"
)

func DeleteFile(ctx context.Context, reqCtx *app.RequestContext) {
  // 从请求中获取待删除文件的路径
  filePath := reqCtx.Query("uri")
  if filePath == "" {
    reqCtx.String(http.StatusBadRequest, "url can not be emtpy")
    return
  }
  if !strings.HasPrefix(filePath, "file") {
    reqCtx.JSON(http.StatusInternalServerError, utils.H{
      "code":    0,
      "message": "url must starts with file",
    })
    return
  }

  // 检查文件是否存在
  if _, err := os.Stat(filePath); os.IsNotExist(err) {
    reqCtx.JSON(http.StatusInternalServerError, utils.H{
      "code":    0,
      "message": "no such file",
    })
    return
  }

  // 尝试删除文件
  if err := os.Remove(filePath); err != nil {
    reqCtx.JSON(http.StatusInternalServerError, utils.H{
      "code":    0,
      "message": err.Error(),
    })
    return
  } else {
    err := DeleteFileAndFramesByUrl(filePath)
    if err != nil {
      reqCtx.JSON(http.StatusInternalServerError, utils.H{
        "code":    0,
        "message": err.Error(),
      })
    }
  }

  // 返回删除成功的响应
  reqCtx.JSON(http.StatusOK, utils.H{
    "code":    200,
    "message": "deleted",
  })
}
