package myutils

import (
  "bytes"
  "crypto/md5"
  "encoding/hex"
  "fmt"
  "io"
  "os"
  "os/exec"
  "path/filepath"
  "strconv"
  "strings"
  "time"
)

// GenerateSnowflakeID 生成唯一的雪花 ID
func GenerateSnowflakeID() string {
  id, err := Flake.NextID()
  if err != nil {
    fmt.Println("Failed to generate snowflake ID:", err)
    return fmt.Sprintf("%d", time.Now().UnixNano())
  }
  return fmt.Sprintf("%d", id)
}

// CalculateFileMD5Path 计算文件路径的 MD5
func CalculateFileMD5Path(filePath string) (string, error) {
  f, err := os.Open(filePath)
  if err != nil {
    return "", err
  }
  defer f.Close()

  hash := md5.New()
  if _, err := io.Copy(hash, f); err != nil {
    return "", err
  }

  return hex.EncodeToString(hash.Sum(nil)), nil
}

// GetVideoDuration 使用 ffprobe 获取视频时长（秒）
func GetVideoDuration(videoPath string) (float64, error) {
  cmd := exec.Command("ffprobe",
    "-v", "error",
    "-show_entries", "format=duration",
    "-of", "default=noprint_wrappers=1:nokey=1",
    videoPath,
  )

  var out bytes.Buffer
  cmd.Stdout = &out
  err := cmd.Run()
  if err != nil {
    return 0, err
  }

  durationStr := strings.TrimSpace(out.String())
  duration, err := strconv.ParseFloat(durationStr, 64)
  if err != nil {
    return 0, err
  }

  return duration, nil
}

// GetAudioDuration 使用 ffprobe 获取音频时长（秒），返回类型为 int8
func GetAudioDuration(audioPath string) (string, error) {
  cmd := exec.Command("ffprobe",
    "-v", "error",
    "-show_entries", "format=duration",
    "-of", "default=noprint_wrappers=1:nokey=1",
    audioPath,
  )

  var out bytes.Buffer
  cmd.Stdout = &out
  err := cmd.Run()
  if err != nil {
    return "", err
  }

  durationStr := strings.TrimSpace(out.String())
  return durationStr, nil
}

// ExtractKeyFrames 使用 FFmpeg 提取视频的关键帧
func ExtractKeyFrames(videoPath string, outputDir string, frameCount int) ([]string, error) {
  // 确保输出目录存在
  if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
    return nil, err
  }

  // 构建输出文件名模式，使用雪花 ID 作为文件名
  outputPattern := filepath.Join(outputDir, "%d.png")

  // 获取视频时长
  duration, err := GetVideoDuration(videoPath)
  if err != nil {
    return nil, err
  }

  // 计算 FPS 以提取所需数量的帧
  fps := float64(frameCount) / duration
  if fps <= 0 {
    fps = 1 // 最小 FPS 为 1
  }

  // 构建 FFmpeg 命令
  cmd := exec.Command("ffmpeg",
    "-i", videoPath,
    "-vf", fmt.Sprintf("fps=%.2f", fps),
    "-q:v", "2", // 设置输出图像质量
    outputPattern,
  )

  // 执行命令并捕获输出
  output, err := cmd.CombinedOutput()
  if err != nil {
    return nil, fmt.Errorf("ffmpeg error: %v, output: %s", err, string(output))
  }

  // 收集生成的帧文件路径
  framePaths := []string{}
  for i := 1; i <= frameCount; i++ {
    framePath := fmt.Sprintf("%d.png", i)
    fullPath := filepath.Join(outputDir, framePath)
    if _, err := os.Stat(fullPath); err == nil {
      framePaths = append(framePaths, fullPath)
    }
  }

  return framePaths, nil
}
