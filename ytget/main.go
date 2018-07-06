package main

import (
  "os"
  "fmt"
  "time"
)

type Youtube struct {
  // for extractor
	public_url         string
  // for downloader
  hidden_url         string
	contentLength      float64
	totalWrittenBytes  float64
	downloadLevel      float64
  // for encoder
  fileStem           string
	title              string
	artist             string
	album              string
	year               string
  // for main
  result             error
}

func main() {

  start := time.Now() // start program run timer

  // get filesystem locations
  home := os.Getenv("HOME")
  outDir := home + "/Music/"
  queueFile := outDir + "/.meta/queue.csv"

  // get list of targets and initialize
  Targets := parseCSV(queueFile, outDir)

  // create buffered channels
  jobs := make(chan *Youtube, len(Targets))
  pipe1 := make(chan *Youtube, len(Targets))
  pipe2 := make(chan *Youtube, len(Targets))
  reply := make(chan error, len(Targets))

  // spool up services
  go extractor(jobs, pipe1, reply)
  go downloader(pipe1, pipe2, reply)
  go encoder(pipe2, reply)

  // feed targets into pipeline
  for _, job := range Targets {
    jobs <- job
  }
  close(jobs)

  // read success state
  for result := range reply {
    fmt.Printf("Error:\t%s\n", result)
  }

  // print program runtime
  fmt.Printf("\nRun Time: %s\n", time.Since(start))
}
