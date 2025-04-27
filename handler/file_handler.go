package handler

import (
	"context"
	"fmt"
	"io"
	"mime"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/litongjava/openfile-server/can"
	"github.com/litongjava/openfile-server/myutils"
)

func GetFile(ctx context.Context, reqCtx *app.RequestContext) {
	filepathParam := reqCtx.Param("filepath")

	// 只对 png/jpg/jpeg/gif 做缩略图匹配
	thumbRegex := regexp.MustCompile(`(?i)^(.+?)_(\d+)x(\d+)\.(png|jpe?g|gif)$`)
	if matches := thumbRegex.FindStringSubmatch(filepathParam); len(matches) == 4 {
		filename := matches[1]
		width, _ := strconv.Atoi(matches[2])
		height, _ := strconv.Atoi(matches[3])
		ext := matches[4]

		original := fmt.Sprintf("%s/%s.%s", can.DEFAULT_FILE_PATH, filename, ext)
		thumbDir := path.Join(can.DEFAULT_FILE_PATH, "thumbnails")
		thumbFile := fmt.Sprintf("%s/%s_%dx%d.%s", thumbDir, filename, width, height, ext)

		if _, err := os.Stat(thumbFile); os.IsNotExist(err) {
			if err := os.MkdirAll(thumbDir, 0755); err != nil {
				reqCtx.String(consts.StatusInternalServerError, "mkdir thumbnails failed: "+err.Error())
				return
			}
			if err := myutils.GenerateThumbnail(original, thumbFile, width, height); err != nil {
				reqCtx.String(consts.StatusInternalServerError, "generate thumbnail failed: "+err.Error())
				return
			}
		}
		serveFileWithRange(reqCtx, thumbFile)
		return
	}

	// 普通文件
	filePath := path.Join(can.DEFAULT_FILE_PATH, filepathParam)
	stat, err := os.Stat(filePath)
	if err != nil {
		reqCtx.String(consts.StatusNotFound, "file not found")
		return
	}
	if stat.IsDir() {
		reqCtx.String(consts.StatusBadRequest, "only support file")
		return
	}
	serveFileWithRange(reqCtx, filePath)
}

func serveFileWithRange(reqCtx *app.RequestContext, filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		reqCtx.String(consts.StatusNotFound, "file not found")
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		reqCtx.String(consts.StatusInternalServerError, "stat file error")
		return
	}
	size := stat.Size()

	// 设置 Content-Type & 标记支持 Range
	mimeType := mime.TypeByExtension(path.Ext(filePath))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	reqCtx.SetContentType(mimeType)
	reqCtx.Header("Accept-Ranges", "bytes")

	// 处理浏览器的 Range 请求
	rangeHdr := string(reqCtx.Request.Header.Peek("Range"))
	if strings.HasPrefix(rangeHdr, "bytes=") {
		parts := strings.Split(strings.TrimPrefix(rangeHdr, "bytes="), "-")
		start, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			reqCtx.Status(consts.StatusRequestedRangeNotSatisfiable)
			return
		}
		var end int64
		if parts[1] != "" {
			end, err = strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				reqCtx.Status(consts.StatusRequestedRangeNotSatisfiable)
				return
			}
		} else {
			end = size - 1
		}
		if start > end || end >= size {
			reqCtx.Status(consts.StatusRequestedRangeNotSatisfiable)
			return
		}

		// Partial Content Response
		reqCtx.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, size))
		reqCtx.Header("Content-Length", strconv.FormatInt(end-start+1, 10))
		reqCtx.Status(consts.StatusPartialContent)

		// 游标移动到 start，然后用 SetBodyStream 只写指定长度
		f.Seek(start, io.SeekStart)
		reqCtx.SetBodyStream(f, int(end-start+1)) // <– 使用 SetBodyStream 而非直接写 Writer :contentReference[oaicite:0]{index=0}
		return
	}

	// 普通 200 OK 全量返回
	reqCtx.Header("Content-Length", strconv.FormatInt(size, 10))
	reqCtx.Status(consts.StatusOK)
	reqCtx.SetBodyStream(f, -1) // <– size<0 表示读到 EOF :contentReference[oaicite:1]{index=1}
}
