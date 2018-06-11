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

func (y *Youtube) createFilename() string {

  // get title of target video
  title := y.StreamList[0]["title"]
  title = strings.Replace(title, " ", "_", -1)
  filePath, _ := filepath.Abs(title)

  // get the file extension of the video
  tmp := y.StreamList[0]["type"]
  re := regexp.MustCompile(`[\w]+\/(\w{3});`)
  ext := re.FindStringSubmatch(tmp)

  // return the absolute path of the new file
  return filePath + "." + ext[1]
}

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

  fmt.Printf("\n")
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
    fmt.Printf("\rDownloading: %2.0f%% Complete\t\t\t", y.downloadLevel)
		y.DownloadPercent <- int64(y.downloadLevel)
	}
	return
}
