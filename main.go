package main

import (
  "flag"
  "github.com/cloudwego/hertz/pkg/app/server"
  "github.com/hertz-contrib/cors"
  "github.com/litongjava/openfile-server/router"
  "strconv"
  "time"
)

func main() {
  port := flag.Int("port", 9000, "server port.")
  flag.Parse()
  addr := ":" + strconv.Itoa(*port)

  h := server.New(server.WithHostPorts("0.0.0.0"+addr), server.WithMaxRequestBodySize(600<<20))

  //跨域中间件
  //app.HandlerFunc
  corsFunction := cors.New(cors.Config{
    AllowAllOrigins:  true,                                     //允许所有origin的请求
    AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"}, //允许的方法
    AllowHeaders:     []string{"Origin"},                       //允许的头部
    ExposeHeaders:    []string{"Content-Length"},               //暴漏的头部信息
    AllowCredentials: true,                                     //允许携带证书
    AllowWildcard:    true,                                     //允许使用通配符匹配
    MaxAge:           12 * time.Hour,                           //请求缓存的最长时间
  })

  h.Use(corsFunction)

  router.RegisterHadlder(h)
  h.Spin()
}
