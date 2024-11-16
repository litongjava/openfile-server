package handler

import (
  "context"
  "fmt"
  "github.com/cloudwego/hertz/pkg/app"
  "github.com/cloudwego/hertz/pkg/common/hlog"
  "github.com/cloudwego/hertz/pkg/common/utils"
  "github.com/cloudwego/hertz/pkg/protocol/consts"
  "github.com/litongjava/openfile-server/can"
  "github.com/litongjava/openfile-server/myutils"
  "net/http"
  "os"
  "path"
  "regexp"
  "strconv"
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
  existingURL, err := GetExistingFileURL(md5Sum)
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
      "code":   200,
      "imgUrl": myutils.GetFullHostURL(reqCtx.URI()),
      "data":   existingURL,
      "md5":    md5Sum,
    })
    return
  }

  reqCtx.JSON(http.StatusOK, utils.H{
    "code": -1,
    "msg":  "not exsits",
  })
}

func GetFile(ctx context.Context, reqCtx *app.RequestContext) {
  filepath := reqCtx.Param("filepath")

  //regex := regexp.MustCompile(`^/(.*)_(\d+)x(\d+)\.(.+)$`)
  regex := regexp.MustCompile(`^(.+)_(\d+)x(\d+)\.(\w+)$`)
  matches := regex.FindStringSubmatch(filepath)
  if len(matches) == 0 {
    reqCtx.File(path.Join(can.DEFAULT_FILE_PATH, filepath))
    return
  }

  // 提取文件名和需要的宽高
  filename := matches[1]
  width, _ := strconv.Atoi(matches[2])
  height, _ := strconv.Atoi(matches[3])
  ext := matches[4]

  // 原始文件路径
  originalFilePath := fmt.Sprintf(can.DEFAULT_FILE_PATH+"/%s.%s", filename, ext)

  // 缩略图文件路径
  thumbnailFilePath := fmt.Sprintf(can.DEFAULT_FILE_PATH+"/thumbnails/%s_%dx%d.%s", filename, width, height, ext)

  // 检查缩略图是否存在
  if _, err := os.Stat(thumbnailFilePath); os.IsNotExist(err) {
    // 如果缩略图不存在，生成缩略图
    if err := myutils.GenerateThumbnail(originalFilePath, thumbnailFilePath, width, height); err != nil {
      reqCtx.JSON(consts.StatusInternalServerError, fmt.Sprintf("Failed to generate thumbnail: %v", err))
      return
    }
  }

  // 返回缩略图
  reqCtx.File(thumbnailFilePath)
}
