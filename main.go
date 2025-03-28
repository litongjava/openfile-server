package main

import (
  "flag"
  "github.com/cloudwego/hertz/pkg/app/server"
  "github.com/cloudwego/hertz/pkg/common/hlog"
  "github.com/hertz-contrib/cors"
  "github.com/litongjava/openfile-server/can"
  "github.com/litongjava/openfile-server/router"
  "io"
  "os"
  "strconv"
  "time"
)

func main() {
  hlog.SetLevel(hlog.LevelDebug)
  f, err := os.Create("app.log")
  if err != nil {
    panic(err)
  }
  defer f.Close()

  // SetOutput sets the output of default logger. By default, it is stderr.
  //hlog.SetOutput(f)
  // if you want to output the log to the file and the stdout at the same time, you can use the following codes
  fileWriter := io.MultiWriter(f, os.Stdout)
  hlog.SetOutput(fileWriter)

  port := flag.Int("port", 9000, "server port.")
  flag.Parse()
  can.OpenDb()
  addr := ":" + strconv.Itoa(*port)

  h := server.New(server.WithHostPorts("0.0.0.0"+addr), server.WithMaxRequestBodySize(600<<20))

  //跨域中间件
  //app.HandlerFunc
  corsFunction := cors.New(cors.Config{
    AllowAllOrigins:  true,                                     // 允许所有 origin 的请求
    AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"}, // 允许的方法
    AllowHeaders:     []string{"Origin", "Content-Type"},       // 允许的头部（添加了 Content-Type）
    ExposeHeaders:    []string{"Content-Length"},               // 暴露的头部信息
    AllowCredentials: true,                                     // 允许携带证书
    AllowWildcard:    true,                                     // 允许使用通配符匹配
    MaxAge:           12 * time.Hour,                           // 请求缓存的最长时间
  })

  h.Use(corsFunction)

  router.RegisterHadlder(h)
  h.Spin()
}
