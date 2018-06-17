package main

import (
  "io"
  "os"
  "fmt"
  "errors"
  "regexp"
  "strings"
  "net/http"
  "path/filepath"
)

func (y *Youtube) createFilename() {
  nospace := regexp.MustCompile(`\s`)
  nospec  := regexp.MustCompile(`[^\w\s]`)

  tmp := fmt.Sprintf("%s %s", y.artist, y.title)
  tmp = strings.ToLower(tmp)
  tmp = nospec.ReplaceAllString(tmp, "")
  tmp = nospace.ReplaceAllString(tmp, "_")

  y.fileStem = tmp
  return
}

func (y *Youtube) videoDLWorker() error {

  // get the content
	resp, err := http.Get(y.tgt_url)
	if err != nil { return err }
  defer resp.Body.Close()


	y.contentLength = float64(resp.ContentLength)  // length of contents

  // check that YouTube returns StatusCode "200"
	if resp.StatusCode != 200 {
		return errors.New("non 200 status code received")
	}

	err = os.MkdirAll(filepath.Dir(y.fileStem), 666)
	if err != nil { return err }

	out, err := os.Create(y.fileStem)
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
    y.downloadPercent <- y.downloadLevel
	}
	return
}
