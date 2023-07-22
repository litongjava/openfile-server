package main

import (
  "github.com/gin-gonic/gin"
  "github.com/litongjava/openfile-server/controller"
  "log"
  "net/http"
  "os"
)

func init() {
  //设置Flats为 日期 时间 微秒 文件名:行号
  log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

func main() {
  log.Println("start...")
  port := os.Getenv("PORT")
  if port == "" {
    log.Fatal("$PORT must be set")
  }

  gin.SetMode(gin.ReleaseMode)
  r := gin.New()
  r.Use(gin.Recovery())

  r.GET("/ping", controller.Ping)
  r.POST("/upload/:username/:repositoryName/*subFolder", controller.Upload)
  r.StaticFS("/s", http.Dir("./storage"))
  r.Run(":" + port)

}
