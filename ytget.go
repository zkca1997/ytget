package main

import (
  "os"
  "log"
  "fmt"
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

func NewYoutube() *Youtube {
	return &Youtube{ DownloadPercent: make(chan int64, 100) }
}

func main() {

  if len(os.Args) != 2 {
    log.Fatalf("ytget-go expects 1 argument")
  }

  y := NewYoutube()

  err := y.DecodeURL( os.Args[1] )
  if err != nil {
    log.Fatalf("Failed decoding URL: %s :%s", os.Args[1], err)
  }

  currentFile := y.createFilename()
  err = y.StartDownload(currentFile)
  if err != nil {
    log.Fatalf("Failed downloading URL: %s", err)
  }

  fmt.Printf("\n\n")

}
