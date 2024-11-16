package handler

import (
  "database/sql"
  "github.com/litongjava/openfile-server/can"
)

func SaveFileInfoToDB(md5Sum, filePath string) error {
  insertSQL := "INSERT INTO open_files(md5,url) VALUES(?,?)"
  _, err := can.Db.Exec(insertSQL, md5Sum, filePath)
  return err
}

func GetExistingFileURL(md5Sum string) (string, error) {
  selectSQL := "SELECT url FROM open_files WHERE md5=?"
  var url string
  err := can.Db.QueryRow(selectSQL, md5Sum).Scan(&url)
  if err == sql.ErrNoRows {
    return "", nil
  }
  return url, err
}
