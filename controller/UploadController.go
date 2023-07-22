package controller

import (
  "fmt"
  "github.com/gin-gonic/gin"
  "log"
  "net/http"
  "path/filepath"
)

func Upload(c *gin.Context) {
  username := c.Param("username")             // 获取url中的username
  repositoryName := c.Param("repositoryName") // 获取url中的repositoryName
  subFolder := c.Param("subFolder")           // 获取url中的filePath

  file, err := c.FormFile("file") // 获取上传的文件
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error": err.Error(),
    })
    return
  }

  filename := file.Filename
  path := filepath.Join("storage", username, repositoryName, subFolder, filename) // 构建保存的文件路径
  log.Println("path:", path)

  if err := c.SaveUploadedFile(file, path); err != nil { // 保存文件
    c.JSON(http.StatusInternalServerError, gin.H{
      "error": err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "message": fmt.Sprintf("'%s' uploaded!", filename),
  })
}
