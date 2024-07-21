package myutils

import (
  "github.com/cloudwego/hertz/pkg/protocol"
)

// getFullHostURL constructs the full host URL including the protocol.
func GetFullHostURL(uri *protocol.URI) string {
  var scheme = string(uri.Scheme())
  var host = string(uri.Host())
  return scheme + "://" + host + "/"
}
