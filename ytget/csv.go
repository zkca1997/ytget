package main

import (
  "os"
  "io"
  "log"
  "bufio"
  "regexp"
  "strings"
  "encoding/csv"
  "path/filepath"
)

func parseCSV(inFile string, directory string) []*Youtube {

  // define regexes for naming conventions
  nospace := regexp.MustCompile(`\s`)
  nospec  := regexp.MustCompile(`[^\w\s]`)

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

    track := strings.ToLower(line[1])
    track = nospec.ReplaceAllString(track, "")
    track = nospace.ReplaceAllString(track, "_")

    artist := strings.ToLower(line[2])
    artist = nospec.ReplaceAllString(artist, "")
    artist = nospace.ReplaceAllString(artist, "_")

    album := strings.ToLower(line[3])
    album = nospec.ReplaceAllString(album, "")
    album = nospace.ReplaceAllString(album, "_")

    Targets = append(Targets, &Youtube{
      public_url: line[0],
      title:      line[1],
      artist:     line[2],
      album:      line[3],
      year:       line[4],
      fileStem:   filepath.Join(directory, artist, album , track),
    })
  }

  return Targets
}
