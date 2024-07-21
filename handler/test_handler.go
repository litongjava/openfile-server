package handler

import (
  "context"
  "github.com/cloudwego/hertz/pkg/app"
  "github.com/cloudwego/hertz/pkg/common/utils"
)

func TestHandler(ctx context.Context, reqCtx *app.RequestContext) {

  reqCtx.JSON(200, utils.H{
    "url":    reqCtx.Request.URI().String(),
    "path":   string(reqCtx.URI().Path()),
    "schema": string(reqCtx.URI().Scheme()),
    "host":   string(reqCtx.URI().Host()),
  })
}
