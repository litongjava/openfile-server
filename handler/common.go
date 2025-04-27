package handler

import (
  "database/sql"
  "fmt"
  "github.com/cloudwego/hertz/pkg/common/hlog"
  "github.com/litongjava/openfile-server/can"
  "github.com/litongjava/openfile-server/myutils"
  "os"
  "path/filepath"
)

func SaveFileInfoToDB(md5Sum, filePath string) error {
  insertSQL := "INSERT INTO open_files(md5,url) VALUES(?,?)"
  _, err := can.Db.Exec(insertSQL, md5Sum, filePath)
  return err
}

func GetMd5ByFiepath(filepath string) (string, error) {
  selectSQL := "SELECT md5 FROM open_files WHERE url=?"
  var md5 string
  err := can.Db.QueryRow(selectSQL, filepath).Scan(&filepath)
  if err == sql.ErrNoRows {
    return "", nil
  }
  return md5, err
}

func SaveAudioFileInfoToDB(md5Sum, filePath, extra string) error {
  insertSQL := "INSERT INTO open_files(md5,url,extra) VALUES(?,?,?)"
  _, err := can.Db.Exec(insertSQL, md5Sum, filePath, extra)
  return err
}

func GetFilepathFromDb(md5Sum string) (string, error) {
  selectSQL := "SELECT url FROM open_files WHERE md5=?"
  var url string
  err := can.Db.QueryRow(selectSQL, md5Sum).Scan(&url)
  if err == sql.ErrNoRows {
    return "", nil
  }
  return url, err
}

func SaveVideoFramesToDB(md5Sum, filePath, frames string) error {
  insertSQL := "INSERT INTO open_file_frames(md5,url,frames) VALUES(?,?,?)"
  _, err := can.Db.Exec(insertSQL, md5Sum, filePath, frames)
  return err
}

func DeleteFileAndFramesByUrl(url string) error {
  // 开启事务
  tx, err := can.Db.Begin()
  if err != nil {
    return err
  }

  // 定义一个用于回滚的辅助函数
  rollback := func() {
    _ = tx.Rollback()
  }

  // 从 open_files 表中删除记录
  deleteFilesSQL := "DELETE FROM open_files WHERE url = ?"
  if _, err := tx.Exec(deleteFilesSQL, url); err != nil {
    rollback()
    return err
  }

  // 从 open_file_frames 表中删除记录
  deleteFramesSQL := "DELETE FROM open_file_frames WHERE url = ?"
  if _, err := tx.Exec(deleteFramesSQL, url); err != nil {
    rollback()
    return err
  }

  // 提交事务
  if err := tx.Commit(); err != nil {
    return err
  }

  return nil
}

func GetVideoFramesFromDb(uri string) (error, string) {
  selectSQL := "SELECT frames FROM open_file_frames WHERE url=?"
  var frames string
  err := can.Db.QueryRow(selectSQL, uri).Scan(&frames)
  if err == sql.ErrNoRows {
    return nil, ""
  }
  return err, frames
}

func QueryAudioLengthFromDb(uri string) (error, string) {
  selectSQL := "SELECT extra FROM open_files WHERE url=?"
  var frames string
  err := can.Db.QueryRow(selectSQL, uri).Scan(&frames)
  if err == sql.ErrNoRows {
    return nil, ""
  }
  return err, frames
}

func ExtraFrames(filePath string, fold string) []string {
  var frames []string
  // 获取视频时长
  duration, err := myutils.GetVideoDuration(filePath)
  if err != nil {
    hlog.Error("Failed to get video duration:", filePath+" ", err)
  } else {
    var frameCount int
    if duration >= 10 {
      frameCount = 10
    } else {
      frameCount = int(duration)
      if frameCount < 1 {
        frameCount = 1
      }
    }

    // 提取关键帧
    frameDir := filepath.Join("file", "frames", fold)
    framePaths, err := myutils.ExtractKeyFrames(filePath, frameDir, frameCount)
    if err != nil {
      hlog.Error("Failed to extract key frames:", err)
    } else {
      for _, framePath := range framePaths {
        // 生成雪花ID作为文件名
        snowflakeID := myutils.GenerateSnowflakeID()
        newFrameFilename := fmt.Sprintf("%s%s", snowflakeID, filepath.Ext(framePath))
        newFramePath := filepath.Join(filepath.Dir(framePath), newFrameFilename)

        // 重命名帧文件，确保使用正斜杠
        err := os.Rename(framePath, newFramePath)
        if err != nil {
          hlog.Error("Failed to rename frame file:", err)
          continue
        }
        var relativeFramePath = filepath.ToSlash(newFramePath)

        // 添加到 frames 列表
        frames = append(frames, relativeFramePath)
      }
    }
  }
  return frames
}
