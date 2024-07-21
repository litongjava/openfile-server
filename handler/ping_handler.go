package handler

import (
  "context"
  "github.com/cloudwego/hertz/pkg/app"
  "github.com/cloudwego/hertz/pkg/common/utils"
  "github.com/cloudwego/hertz/pkg/protocol/consts"
)

func PingHandler(ctx context.Context, reqCtx *app.RequestContext) {
  reqCtx.JSON(consts.StatusOK, utils.H{"message": "pong"})
}
