package router

import (
  "context"
  "fmt"
  "github.com/cloudwego/hertz/pkg/app"
  "github.com/cloudwego/hertz/pkg/app/server"
  "github.com/cloudwego/hertz/pkg/protocol/consts"
  "github.com/litongjava/openfile-server/can"
  "github.com/litongjava/openfile-server/handler"
  "github.com/litongjava/openfile-server/myutils"
  "os"
  "path"
  "regexp"
  "strconv"
  "time"
)

func RegisterHadlder(h *server.Hertz) {
  h.GET("/ping", handler.PingHandler)
  h.GET("/test", handler.TestHandler)

  h.POST("/upload/:username/:repositoryName/*subFolder", handler.UploadHandler)

  h.POST("/uploadImg", func(c context.Context, ctx *app.RequestContext) {
    handler.Upload(ctx, can.DEFAULT_FILE_PATH+"/image/")
  })
  h.POST("/uploadVideo", func(c context.Context, ctx *app.RequestContext) {
    handler.Upload(ctx, can.DEFAULT_FILE_PATH+"/video/")
  })
  h.POST("/uploadMp3", func(c context.Context, ctx *app.RequestContext) {
    handler.Upload(ctx, can.DEFAULT_FILE_PATH+"/audio/")
  })
  h.POST("/uploadDoc", func(c context.Context, ctx *app.RequestContext) {
    handler.Upload(ctx, can.DEFAULT_FILE_PATH+"/doc/")
  })

  h.POST("/upload", func(c context.Context, ctx *app.RequestContext) {

    handler.Upload(ctx, can.DEFAULT_FILE_PATH+"/")
  })

  h.GET("/url", handler.GetUrl)
  h.GET("/file/*filepath", func(ctx context.Context, c *app.RequestContext) {
    filepath := c.Param("filepath")

    //regex := regexp.MustCompile(`^/(.*)_(\d+)x(\d+)\.(.+)$`)
    regex := regexp.MustCompile(`^(.+)_(\d+)x(\d+)\.(\w+)$`)
    matches := regex.FindStringSubmatch(filepath)
    if len(matches) == 0 {
      c.File(path.Join(can.DEFAULT_FILE_PATH, filepath))
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
        c.JSON(consts.StatusInternalServerError, fmt.Sprintf("Failed to generate thumbnail: %v", err))
        return
      }
    }

    // 返回缩略图
    c.File(thumbnailFilePath)
  })

  //h.StaticFS("/file", "./file")
  //h.StaticFS("/file", &app.FS{Root: ""})
  //h.StaticFS("/file", &app.FS{
  //  Root: "",
  //  PathNotFound: func(_ context.Context, ctx *app.RequestContext) {
  //    ctx.JSON(consts.StatusNotFound, "The requested resource does not exist")
  //  },
  //  CacheDuration:   time.Second * 5,
  //  AcceptByteRange: true,
  //})
  h.StaticFS("/s", &app.FS{
    Root: "",
    PathNotFound: func(_ context.Context, ctx *app.RequestContext) {
      ctx.JSON(consts.StatusNotFound, "The requested resource does not exist")
    },
    CacheDuration:      time.Second * 5,
    AcceptByteRange:    true,
    GenerateIndexPages: true,
  })
}
