package main

import (
  "os"
  "io"
  "fmt"
  "regexp"
  "strings"
  "net/http"
  "path/filepath"
)

func downloader(in <-chan *Youtube, out chan<- *Youtube, fail chan<- error) {

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

func (y *Youtube) downloadWorker(out chan<- *Youtube, fail chan<- error, done chan<- bool) {

  // get the content
	resp, err := http.Get(y.hidden_url)
	if err != nil {
    fmt.Println(y.hidden_url)
    fail <- fmt.Errorf("job [%s by %s]: %s", y.title, y.artist, err)
    return
  }
  defer resp.Body.Close()

	y.contentLength = float64(resp.ContentLength)  // length of contents

  // check that YouTube returns StatusCode "200"
	if resp.StatusCode != 200 {
		fail <- fmt.Errorf("job [%s by %s]: non 200 status returned", y.title, y.artist)
    return
	}

  y.createFilename()
	err = os.MkdirAll(filepath.Dir(y.fileStem), 666)
	if err != nil {
    fail <- fmt.Errorf("job [%s by %s]: %s", y.title, y.artist, err)
    return
  }

	output, err := os.Create(y.fileStem)
	if err != nil {
    fail <- fmt.Errorf("job [%s by %s]: %s", y.title, y.artist, err)
    return
  }

	mw := io.MultiWriter(output, y)
	_, err = io.Copy(mw, resp.Body)
  if err != nil {
    fail <- fmt.Errorf("job [%s by %s]: %s", y.title, y.artist, err)
    return
  }

  fmt.Printf("Downloaded:\t%s by %s\n", y.title, y.artist)
  out <- y
  done <- true
	return
}

func (y *Youtube) Write(p []byte) (n int, err error) {
	n = len(p)
	y.totalWrittenBytes = y.totalWrittenBytes + float64(n)
	currentPercent := ((y.totalWrittenBytes / y.contentLength) * 100)
	if (y.downloadLevel <= currentPercent) && (y.downloadLevel < 100) {
		y.downloadLevel++
	}
	return
}

func (y *Youtube) createFilename() {
  nospace := regexp.MustCompile(`\s`)
  nospec  := regexp.MustCompile(`[^\w\s]`)

  tmp := fmt.Sprintf("%s %s", y.artist, y.title)
  tmp = strings.ToLower(tmp)
  tmp = nospec.ReplaceAllString(tmp, "")
  tmp = nospace.ReplaceAllString(tmp, "_")

  y.fileStem = y.directory + tmp
  return
}
