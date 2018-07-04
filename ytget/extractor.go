package main

import (
  "log"
  "fmt"
  "time"
  "os/exec"
  "net/http"
  "io/ioutil"
)

func extractor(in <-chan *Youtube, out chan<- *Youtube, fail chan<- error) {

  extractorService := exec.Command("extractorService.py")
  err := extractorService.Start()
  if err != nil { log.Fatal(err) }

  // test server every 10 miliseconds for up to 5 seconds
  waittime := 10
  testurl := "http://localhost:8081?input=test"
  for keepalive := 5000; true; keepalive -= waittime {
    if keepalive < 0 {
      log.Fatalf("extractorService timeout")
      if extractorService.Process.Kill() != nil {
        log.Fatal("failed to bring down 'extractorService'")
      }
    }
    resp, err := http.Get(testurl)
    if err == nil {
      resp.Body.Close()
      if resp.StatusCode == 400 { break }
    }
    time.Sleep(time.Duration(waittime) * time.Millisecond)
  }

  // extract hidden urls from jobs as they arrive on the buffer
  for job := range in {
    err = job.extractorRequest()
    if err != nil { fail <- err
    } else {
      out <- job
      fmt.Printf("Extracted:\t%s by %s\n", job.title, job.artist)
    }
  }
  close(out)

  // done, bring down extractorService
  if extractorService.Process.Kill() != nil {
    log.Fatal("failed to bring down 'extractorService'")
  }
}

func (y *Youtube) extractorRequest() error {

  // query extractor server for hidden url
  req_url := "http://localhost:8081?input=" + y.public_url
  resp, err := http.Get(req_url)
  if err != nil {
    return fmt.Errorf("job [%s by %s]: %s", y.title, y.artist, err)
  }
  defer resp.Body.Close()

  // check status code
  if resp.StatusCode != 200 {
    return fmt.Errorf("job [%s by %s]: Server Reponse Status not 200", y.title, y.artist)
  }

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return fmt.Errorf("job [%s by %s]: %s", y.title, y.artist, err)
  }

  y.hidden_url = string(body)
  return nil
}
