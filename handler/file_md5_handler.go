package handler

import (
  "context"
  "fmt"
  "github.com/cloudwego/hertz/pkg/app"
  "github.com/cloudwego/hertz/pkg/common/utils"
  "github.com/cloudwego/hertz/pkg/protocol/consts"
  "net/http"
  "os"
)

func QueryFileByMD5Handler(ctx context.Context, reqCtx *app.RequestContext) {
  md5Sum := reqCtx.Param("md5")
  if md5Sum == "" {
    reqCtx.JSON(consts.StatusBadRequest, utils.H{
      "code":    "0",
      "message": "md5 can not emtpy",
    })
    return
  }

  // 根据 md5 查询数据库中的文件路径
  filePath, err := GetFilepathFromDb(md5Sum)
  if err != nil {
    reqCtx.JSON(consts.StatusNotFound, utils.H{
      "code":    0,
      "message": fmt.Sprintf("not found from database：%v", err),
    })
    return
  }

  if filePath == "" {
    reqCtx.JSON(http.StatusNotFound, utils.H{
      "code":    0,
      "message": "no such file",
    })
    return
  }

  // 检查磁盘上是否存在该文件
  if _, err := os.Stat(filePath); os.IsNotExist(err) {
    reqCtx.JSON(http.StatusOK, utils.H{
      "code":    0,
      "message": "file not found",
    })
    return
  } else if err != nil {
    reqCtx.JSON(consts.StatusInternalServerError, utils.H{
      "code":    0,
      "message": fmt.Sprintf("failed to check file status：%v", err),
    })
    return
  }

  // 文件存在，返回成功信息
  reqCtx.JSON(http.StatusOK, utils.H{
    "code": 200,
    "url":  filePath,
  })
}
