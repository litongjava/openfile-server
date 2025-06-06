package handler

import (
  "context"
  "github.com/cloudwego/hertz/pkg/app"
  "github.com/cloudwego/hertz/pkg/common/hlog"
  "github.com/cloudwego/hertz/pkg/common/utils"
  "github.com/litongjava/openfile-server/myutils"
  "net/http"
  "os"
)

func GetUrl(ctx context.Context, reqCtx *app.RequestContext) {
  md5Sum, hasMd5Sum := reqCtx.GetQuery("md5")
  //= reqCtx.GetPostForm("md5")
  if !hasMd5Sum {
    reqCtx.JSON(http.StatusOK, utils.H{
      "code": -1,
      "msgs": "md5 can not be empty",
    })
  }
  hlog.Info("md5:", md5Sum)
  // Check if file already exists in DB
  existingURL, err := GetFilepathFromDb(md5Sum)
  if err == nil && existingURL != "" {
    _, err := os.Stat(existingURL)
    if err != nil {
      reqCtx.JSON(http.StatusOK, utils.H{
        "code": -1,
        "msg":  err.Error(),
      })
      return
    }

    reqCtx.JSON(http.StatusOK, utils.H{
      "code": 200,
      "url":  myutils.GetFullHostURL(reqCtx.URI()),
      "data": existingURL,
      "md5":  md5Sum,
    })
    return
  }

  reqCtx.JSON(http.StatusOK, utils.H{
    "code": -1,
    "msg":  "not exsits",
  })
}
