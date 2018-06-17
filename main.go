package main

import (
  "os"
  "fmt"
  "time"
  "strings"
  "strconv"
  "os/exec"
)

type Youtube struct {
  tgt_url           string
  fileStem          string
	contentLength     float64
	totalWrittenBytes float64
	downloadLevel     float64
  downloadPercent   chan float64
  url               string
  title             string
  artist            string
  album             string
  year              string
}

func worker(id int, jobs <-chan Youtube, reply chan<- error, done chan<- bool) {

  for job := range jobs {
    reply <- DownloadURL(&job)
    fmt.Printf("Channel %d: %s by %s\n", id, job.title, job.artist)
  }
  done <- true
}

func displayAll(active []Youtube) {

  fmt.Println("\033[2J")
  fmt.Print("\nDownload Status:\n----------------\n")

  // print out active status
  for _, y := range active {
    fmt.Printf("[%3.0f%%] Downloading: %s by %s\n",
      y.downloadLevel, y.title, y.artist)
  }
}

func DownloadURL(y *Youtube) error {

  // base name of target file
  y.createFilename()

  // decode the url using Youtube-DL plugin
  err := y.DecodeURL()
  if err != nil {
    return fmt.Errorf("Failed decoding %s by %s: %s", y.title, y.artist, err)
  }

  // begin the download
  err = y.videoDLWorker()
  if err != nil {
    return fmt.Errorf("Failed downloading %s by %s: %s", y.title, y.artist, err)
  }

  // convert the file to raw WAV audio file
  wavFile, err := y.toWAV()
  if err != nil {
    return fmt.Errorf("Failed to decompress %s by %s: %s", y.title, y.artist, err)
  }

  // re-encode WAV file to opus
  err = y.toOPUS(wavFile)
  if err != nil {
    return fmt.Errorf("Failed to encode %s by %s: %s", y.title, y.artist, err)
  }

  // clean up the temporary files
  err = cleanFile(y.fileStem, wavFile)
  if err != nil {
    return fmt.Errorf("Failed to clean residual files for %s by %s: %s", y.title, y.artist, err)
  }

  return nil
}

func (y *Youtube) DecodeURL() error {
  cmd := exec.Command("python3", "ytdl.py", y.url)

  out, err := cmd.Output()
  if err != nil {
    return err
  }

  y.tgt_url = strings.TrimSpace(string(out))
  return nil
}

func main() {

  start := time.Now()

  // check command line arguments
  if len(os.Args) != 3 {
    fmt.Println("USAGE: ytget [path to csv file] [# workers]")
    return
  }

  Targets := parseCSV(os.Args[1])
  workers, err := strconv.Atoi(os.Args[2])  // number of worker to run
  if err != nil {
    fmt.Println("Invalid number of workers entered")
    return
  }

  // limit the number of workers that can be spawned
  if workers < 1 || workers > 10 {
    fmt.Println("Must have between 1 and 10 workers")
    return
  }

  // make sure more workers are not started than jobs
  if len(Targets) < workers { workers = len(Targets) }

  // create buffered channels
  jobs := make(chan Youtube, len(Targets))
  reply := make(chan error, len(Targets))

  // start workers
  var done []chan bool
  for w := 1; w <= workers; w++ {
    trigger := make(chan bool, 1)
    go worker(w, jobs, reply, trigger)
    done = append(done, trigger)
  }

  // stack jobs onto the work queue
  for _, job := range Targets {
    jobs <- job
  }
  close(jobs)

  // when each worker has finished, collect errors
  for _, msg := range done {
    <-msg
  }
  close(reply)

  // print out errors as they arrive
  for result := range reply {
    if result != nil {
      fmt.Println(result)
    }
  }

  // program run-time
  fmt.Printf("\nRun Time: %s\n", time.Since(start))
}
