package myutils

import (
  "github.com/sony/sonyflake"
  "time"
)

// Initialize Sonyflake for generating unique IDs
var Flake *sonyflake.Sonyflake

func init() {
  Flake = sonyflake.NewSonyflake(sonyflake.Settings{
    StartTime: time.Now(),
  })
  if Flake == nil {
    panic("Failed to initialize sonyflake")
  }
}
