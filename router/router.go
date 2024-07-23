package router

import (
  "context"
  "github.com/cloudwego/hertz/pkg/app"
  "github.com/cloudwego/hertz/pkg/app/server"
  "github.com/cloudwego/hertz/pkg/protocol/consts"
  "github.com/litongjava/openfile-server/handler"
  "time"
)

func RegisterHadlder(h *server.Hertz) {
  h.GET("/ping", handler.PingHandler)
  h.GET("/test", handler.TestHandler)

  h.POST("/upload/:username/:repositoryName/*subFolder", handler.UploadHandler)

  h.POST("/uploadImg", func(c context.Context, ctx *app.RequestContext) {
    handler.Upload(ctx, "file/image/")
  })
  h.POST("/uploadVideo", func(c context.Context, ctx *app.RequestContext) {
    handler.Upload(ctx, "file/video/")
  })
  h.POST("/uploadMp3", func(c context.Context, ctx *app.RequestContext) {
    handler.Upload(ctx, "file/video/")
  })
  h.POST("/uploadDoc", func(c context.Context, ctx *app.RequestContext) {
    handler.Upload(ctx, "file/doc/")
  })

  h.POST("/upload", func(c context.Context, ctx *app.RequestContext) {
    handler.Upload(ctx, "file/")
  })

  h.GET("/url", handler.GetUrl)

  //h.StaticFS("/file", "./file")
  //h.StaticFS("/file", &app.FS{Root: ""})
  h.StaticFS("/file", &app.FS{
    Root: "",
    PathNotFound: func(_ context.Context, ctx *app.RequestContext) {
      ctx.JSON(consts.StatusNotFound, "The requested resource does not exist")
    },
    CacheDuration:      time.Second * 5,
    AcceptByteRange:    true,
    GenerateIndexPages: true,
  })
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
