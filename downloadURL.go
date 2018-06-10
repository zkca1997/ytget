package main

import (
  "io"
  "os"
  "errors"
  "net/http"
  "path/filepath"
)

func (y *Youtube) StartDownload(destFile string) error {

  //download highest resolution on [0]
	targetStream := y.StreamList[0]
	url := targetStream["url"] + "&signature=" + targetStream["sig"]

	err := y.videoDLWorker(destFile, url)
	return err
}

func (y *Youtube) videoDLWorker(destFile string, target string) error {

  // get the content
	resp, err := http.Get(target)
  defer resp.Body.Close()
	if err != nil { return err }

	y.contentLength = float64(resp.ContentLength)  // length of contents

  // check that YouTube returns StatusCode "200"
	if resp.StatusCode != 200 {
		return errors.New("non 200 status code received")
	}

	err = os.MkdirAll(filepath.Dir(destFile), 666)
	if err != nil { return err }

	out, err := os.Create(destFile)
	if err != nil { return err }

	mw := io.MultiWriter(out, y)
	_, err = io.Copy(mw, resp.Body)
	if err != nil { return err }

	return nil
}

func (y *Youtube) Write(p []byte) (n int, err error) {
	n = len(p)
	y.totalWrittenBytes = y.totalWrittenBytes + float64(n)
	currentPercent := ((y.totalWrittenBytes / y.contentLength) * 100)
	if (y.downloadLevel <= currentPercent) && (y.downloadLevel < 100) {
		y.downloadLevel++
		y.DownloadPercent <- int64(y.downloadLevel)
	}
	return
}
