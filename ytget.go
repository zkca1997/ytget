package main

import (
  "log"
  "regexp"
  "path/filepath"
)

type stream map[string]string

type Youtube struct {
	StreamList        []stream
	VideoID           string
	videoInfo         string
	DownloadPercent   chan int64
	contentLength     float64
	totalWrittenBytes float64
	downloadLevel     float64
}

func NewYoutube(debug bool) *Youtube {
	return &Youtube{ DownloadPercent: make(chan int64, 100) }
}

func (y *Youtube) createFilename(name string) string {
  filePath, _ := filepath.Abs(name)

  tmp := y.StreamList[0]["type"]
  re := regexp.MustCompile(`[\w]+\/(\w{3});`)
  ext := re.FindStringSubmatch(tmp)

  return filePath + "." + ext[1]
}

func main() {
  y := NewYoutube(true)

  err := y.DecodeURL("https://www.youtube.com/watch?v=-bEQKyyZQmM")
  if err != nil {
    log.Fatalf("Failed decoding URL: %s", err)
  }

  currentFile := y.createFilename("test")
  err = y.StartDownload(currentFile)
  if err != nil {
    log.Fatalf("Failed downloading URL: %s", err)
  }

}
