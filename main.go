package main

import (
  "fmt"
)

type stream map[string]string

type Youtube struct {
	StreamList        []stream
  url               string
	VideoID           string
	videoInfo         string
	DownloadPercent   chan int
	contentLength     float64
	totalWrittenBytes float64
	downloadLevel     float64
  filename          string
  title             string
  artist            string
  album             string
  year              string
}

func worker(id int, jobs <-chan *Youtube, results chan<- string) {
  for job := range jobs {
    fmt.Printf("%d got job %s\n", id, job.title)
    return
    //results <- DownloadURL(job)
  }
}

func DownloadURL(y *Youtube) string {

  // decode the URL and fetch video metadata
  err := y.DecodeURL( y.url )
  if err != nil {
    return fmt.Sprintf("Failed decoding %s by %s: %s\n", y.title, y.artist, err)
  }

  // determine what filetype is being downloaded and set target file
  fileStem, ext := y.createFilename()
  downloadFile := fileStem + ext

  // begin the download
  err = y.StartDownload(downloadFile)
  if err != nil {
    return fmt.Sprintf("Failed downloading %s by %s: %s", y.title, y.artist, err)
  }

  // convert the file to raw WAV audio file
  wavFile, err := y.toWAV(downloadFile)
  if err != nil {
    return fmt.Sprintf("Failed to decompress %s by %s: %s", y.title, y.artist, err)
  }
  // re-encode WAV file to opus
  err = y.toOPUS(wavFile)
  if err != nil {
    return fmt.Sprintf("Failed to encode %s by %s: %s", y.title, y.artist, err)
  }

  // clean up the temporary files
  err = cleanFile(downloadFile, wavFile)
  if err != nil {
    return fmt.Sprintf("Failed to clean residual files for %s by %s: %s", y.title, y.artist, err)
  }

  return fmt.Sprintf("Download Completed: %s by %s", y.title, y.artist)
}

func main() {

  // these will eventually be command like argument parsing
  Targets := parseCSV("test.csv")

  workers := 2  // number of worker to run

  // create buffered channels
  fmt.Println(len(Targets))
  jobs := make(chan *Youtube, len(Targets))
  results := make(chan string, len(Targets))

  // start workers
  for w := 1; w <= workers; w++ {
    go worker(w, jobs, results)
  }

  // stack jobs onto the work queue
  for _, job := range Targets {
    fmt.Printf("pushing job: %s to queue\n", job.title)
    jobs <- &job
  }
  close(jobs)

  // print out return status
/*  for range Targets {
    message := <-results
    fmt.Println(message)
  }
*/
}
