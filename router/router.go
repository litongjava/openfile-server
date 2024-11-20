package router

import (
  "context"
  "github.com/cloudwego/hertz/pkg/app"
  "github.com/cloudwego/hertz/pkg/app/server"
  "github.com/cloudwego/hertz/pkg/protocol/consts"
  "github.com/litongjava/openfile-server/can"
  "github.com/litongjava/openfile-server/handler"
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
    handler.UploadVideo(ctx, can.DEFAULT_FILE_PATH+"/video/")
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
  h.GET("/video/frames", handler.VideoFrames)
  h.GET("/file/*filepath", handler.GetFile)

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
