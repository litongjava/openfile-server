package can

import (
  "github.com/cloudwego/hertz/pkg/common/hlog"
  "github.com/jmoiron/sqlx"
  "github.com/litongjava/openfile-server/config"
  _ "github.com/mattn/go-sqlite3"
  "log"
  "strings"
)

func OpenDb() (err error) {
  dsn := "file:openfile-server.db?cache=shared&mode=rwc" // Example SQLite DSN
  // Connect to the database using sqlx with the SQLite driver
  db, err := sqlx.Connect("sqlite3", dsn)
  if err != nil {
    hlog.Fatal("connect DB failed, err:%v\n", err)
    return
  }
  // For SQLite, these settings are less relevant but setting them doesn't harm
  db.SetMaxOpenConns(20)
  db.SetMaxIdleConns(10)

  Db = db
  createLiteTable(db, "sqlite_init.sql")
  return
}

func createLiteTable(sqlDb *sqlx.DB, createSql string) {
  if createSql == "" {
    hlog.Info("skip create sql")
    return
  }
  //读取createSql的值
  bytes, err := config.InitSql.ReadFile(createSql)
  if err != nil {
    hlog.Fatal("An error occurred while reading the file:", err)
  }

  sql := string(bytes)
  sqlStatements := strings.Split(sql, ";")

  for _, createTableSql := range sqlStatements {
    // 移除字符串两边的空格
    createTableSql = strings.TrimSpace(createTableSql)
    // 检查是否为空字符串
    if createTableSql == "" {
      continue
    }

    _, err := Db.Exec(createTableSql)
    if err != nil {
      log.Println("Error in creating table:", err)
    } else {
      log.Println("SQL execution finished:", createTableSql)
    }
  }
}
