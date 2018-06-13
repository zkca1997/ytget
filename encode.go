package main

import (
  "os"
  "fmt"
  "os/exec"
  "path/filepath"
)

func (y *Youtube) toWAV(downloadFile string) (string, error) {

  wavFile := downloadFile + ".wav"

  cmd := exec.Command("ffmpeg", "-y", "-i", downloadFile, wavFile)
  out, err := cmd.CombinedOutput()
  if err != nil {
    fmt.Println(out)
    return wavFile, err
  }

  return wavFile, nil
}

func (y *Youtube) toOPUS(wavFile string) error {

  //outFile := fmt.Sprintf("%s_%s.opus", y.title, y.artist)
  outFile := wavFile + ".opus"

  cmd := exec.Command("opusenc", "--title", y.title, "--artist", y.artist,
    "--album", y.album, "--date", y.year, wavFile, outFile)

  out, err := cmd.CombinedOutput()
  if err != nil {
    fmt.Printf("%s\n", out)
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
