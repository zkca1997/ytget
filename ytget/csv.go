package main

import (
  "os"
  "io"
  "log"
  "bufio"
  "encoding/csv"
)

func parseCSV(inFile string, directory string) []*Youtube {

  csvFile, err := os.Open(inFile)
  if err != nil {
    log.Fatalf("Error opening CSV file: %s\n", err)
  }

  reader := csv.NewReader(bufio.NewReader(csvFile))
  var Targets []*Youtube

  for {
    line, error := reader.Read()
    if error == io.EOF {
        break
    } else if error != nil {
        log.Fatal(error)
    }

    Targets = append(Targets, &Youtube{
      public_url: line[0],
      title:      line[1],
      artist:     line[2],
      album:      line[3],
      year:       line[4],
      directory:  directory,
    })
  }

  return Targets
}
