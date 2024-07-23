package config

import "embed"

//go:embed sqlite_init.sql
var InitSql embed.FS
