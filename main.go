package main

import (
  "github.com/gin-gonic/gin"
  "github.com/litongjava/openfile-server/controller"
  "log"
  "net/http"
  "os"
)

func init() {
  log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

func main() {
  port := os.Getenv("PORT")
  if port == "" {
    port = "80"
  }

  log.Println("start")

  gin.SetMode(gin.ReleaseMode)
  r := gin.New()
  r.Use(gin.Recovery())

  r.GET("/ping", controller.Ping)
  r.POST("/upload/:username/:repositoryName/*subFolder", controller.Upload)
  r.POST("/u/:username/:repositoryName/*subFolder", controller.Upload)
  r.StaticFS("/s", http.Dir("./storage"))
  r.Run(":" + port)

}
