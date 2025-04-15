package handler

import (
  "context"
  "fmt"
  "os"
  "path"
  "regexp"
  "strconv"

  "github.com/cloudwego/hertz/pkg/app"
  "github.com/cloudwego/hertz/pkg/protocol/consts"
  "github.com/litongjava/openfile-server/can"
  "github.com/litongjava/openfile-server/myutils"
)

func GetFile(ctx context.Context, reqCtx *app.RequestContext) {
  filepathParam := reqCtx.Param("filepath")

  // 使用正则匹配缩略图请求格式：filename_widthxheight.ext
  regex := regexp.MustCompile(`^(.+)_(\d+)x(\d+)\.(\w+)$`)
  matches := regex.FindStringSubmatch(filepathParam)
  if len(matches) == 0 {
    file := path.Join(can.DEFAULT_FILE_PATH, filepathParam)
    if stat, err := os.Stat(file); err == nil && stat.IsDir() {
      reqCtx.String(consts.StatusBadRequest, "only support file")
      return
    }
    reqCtx.File(file)
    return
  }

  // 如果匹配成功，则按缩略图逻辑处理
  // 提取文件名和需要的宽高
  filename := matches[1]
  width, _ := strconv.Atoi(matches[2])
  height, _ := strconv.Atoi(matches[3])
  ext := matches[4]

  originalFilePath := fmt.Sprintf(can.DEFAULT_FILE_PATH+"/%s.%s", filename, ext)

  thumbnailFilePath := fmt.Sprintf(can.DEFAULT_FILE_PATH+"/thumbnails/%s_%dx%d.%s", filename, width, height, ext)

  if _, err := os.Stat(thumbnailFilePath); os.IsNotExist(err) {
    // 如果缩略图不存在，生成缩略图
    if err := myutils.GenerateThumbnail(originalFilePath, thumbnailFilePath, width, height); err != nil {
      reqCtx.JSON(consts.StatusInternalServerError, fmt.Sprintf("Failed to generate thumbnail: %v", err))
      return
    }
  }

  reqCtx.File(thumbnailFilePath)
}
