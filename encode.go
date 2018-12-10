package main

import (
  "os"
  "fmt"
  "os/exec"
  "path/filepath"
  "runtime"
)

func encoder(in <-chan *track, fail chan<- error) {

  done := make(chan bool)

  /* spool up encoder workers equal to the number of system CPUs (this takes
   * maximum advantage of CPU resources provided nothing else is running) */
  workers := runtime.NumCPU()
  for w := 0; w < workers; w++ {
    go encoderWorker(in, fail, done)
  }

  // wait for each worker to report "done"
  for w := 0; w < workers; w++ { <-done }

  // signal main the workflow is over and exit
  close(fail)
  return
}

func encoderWorker(jobs <-chan *track, fail chan<- error, done chan<- bool) {

  for job := range jobs {
    // convert the file to raw WAV audio file
    wavFile, err := job.toWAV()
    if err != nil {
      fail <- fmt.Errorf("Failed to decompress %s: %s", job.humanName(), err)
      continue
    }

    // re-encode WAV file to opus
    err = job.toOPUS(wavFile)
    if err != nil {
      fail <- fmt.Errorf("Failed to encode %s: %s", job.humanName(), err)
      continue
    }

    // clean up the temporary files
    err = cleanFile(job.path, wavFile)
    if err != nil {
      fail <- fmt.Errorf("Failed to clean residual files for %s: %s", job.humanName(), err)
      continue
    }

    fmt.Printf("Encoded:\t%s\n", job.humanName())
  }

  done <- true
  return
}

func (y *track) toWAV() (string, error) {

  wavFile := y.path + ".wav"

  cmd := exec.Command("ffmpeg", "-y", "-i", y.path, wavFile)
  out, err := cmd.CombinedOutput()
  if err != nil {
    fmt.Println(out)
    return wavFile, err
  }

  return wavFile, nil
}

func (y *track) toOPUS(wavFile string) error {

  outFile := y.path + ".ogg"

  cmd := exec.Command("opusenc", "--title", y.title, "--artist", y.artist,
    "--album", y.album, "--date", y.year, wavFile, outFile)

  out, err := cmd.CombinedOutput()
  if err != nil {
    fmt.Println(out)
    return err
  }

  return nil
}

func cleanFile(downloadFile string, wavFile string) error {
  downloadPath, _ := filepath.Abs(downloadFile)
  wavPath, _ := filepath.Abs(wavFile)

  err := os.Remove(downloadPath)
  if err != nil {
    return err
  }

  err = os.Remove(wavPath)
  if err != nil {
    return err
  }

  return nil
}
