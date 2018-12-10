package main

import (
  "os"
  "io"
  "fmt"
  "net/http"
  "path/filepath"
)

func downloader(in <-chan *track, out chan<- *track, fail chan<- error) {

  done := make(chan bool)
  var count int

  for job := range in {
    go job.downloadWorker(out, fail, done)
    count += 1
  }

  // wait for each worker to report "done"
  for i := 0; i < count; i++ { <-done }

  // spin down the downloader
  close(out)
  return
}

func (y *track) downloadWorker(out chan<- *track, fail chan<- error, done chan<- bool) {

  // get the content
	resp, err := http.Get(y.hidden_url)
	if err != nil {
    fmt.Println(y.hidden_url)
    fail <- fmt.Errorf("job [%s]: %s", y.humanName(), err)
    return
  }
  defer resp.Body.Close()

	y.contentLength = float64(resp.ContentLength)  // length of contents

  // check that YouTube returns StatusCode "200"
	if resp.StatusCode != 200 {
		fail <- fmt.Errorf("job [%s]: non 200 status returned", y.humanName())
    return
	}

  // create a working directory for download
	err = os.MkdirAll(filepath.Dir(y.path), 0755)
	if err != nil {
    fail <- fmt.Errorf("job [%s]: %s", y.humanName(), err)
    return
  }

	output, err := os.Create(y.path)
	if err != nil {
    fail <- fmt.Errorf("job [%s]: %s", y.humanName(), err)
    return
  }

	mw := io.MultiWriter(output, y)
	_, err = io.Copy(mw, resp.Body)
  if err != nil {
    fail <- fmt.Errorf("job [%s]: %s", y.humanName(), err)
    return
  }

  fmt.Printf("Downloaded:\t%s\n", y.humanName())
  out <- y
  done <- true
	return
}

func (y *track) Write(p []byte) (n int, err error) {
	n = len(p)
	y.totalWrittenBytes = y.totalWrittenBytes + float64(n)
	currentPercent := ((y.totalWrittenBytes / y.contentLength) * 100)
	if (y.downloadLevel <= currentPercent) && (y.downloadLevel < 100) {
		y.downloadLevel++
	}
	return
}
