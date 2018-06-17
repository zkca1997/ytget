package main

import (
  "os"
  "io"
  "fmt"
  "log"
  "bufio"
  "regexp"
  "strings"
  "io/ioutil"
  "encoding/csv"
)

type track struct {
  url       string
  title     string
  artist    string
  album     string
  year      string
  filename  string
}

func main() {

  music_dir := "/home/tkirk/Music/"
  have := getCurrentTracks(music_dir)

  man_file := "/home/tkirk/Music/.meta/manifest.csv"
  need := missingManifest(man_file, have)

  fmt.Printf("\nAdding to Queue:\n----------------\n")
  for _, entry := range need {
    fmt.Printf("%s by %s\n", entry.title, entry.artist)
  }
  fmt.Println()

  queue_file := "/home/tkirk/Music/.meta/queue.csv"
  writeQueue(need, queue_file)
}

func writeQueue(tracks []track, filename string) {
  file, err := os.Create(filename)
  if err != nil { log.Fatal(err) }
  defer file.Close()

  writer := csv.NewWriter(file)
  defer writer.Flush()

  for _, entry := range tracks {
    var data = []string {
      entry.url,
      entry.title,
      entry.artist,
      entry.album,
      entry.year }
    err := writer.Write(data)
    if err != nil { log.Fatal(err) }
  }
}

func getCurrentTracks(music_path string) []string {

  files, err := ioutil.ReadDir(music_path)
  if err != nil {
      log.Fatal(err)
  }

  var tracks []string
  for _, f := range files {
    tracks = append(tracks, f.Name())
  }

  return tracks
}

func missingManifest(filename string, have[]string) []track {

  // open file
  manifestFile, err := os.Open(filename)
  if err != nil {
    log.Fatal("Error opening %s: %s\n", filename, err)
  }

  // open CSV parser
  reader := csv.NewReader(bufio.NewReader(manifestFile))
  var tracks []track

  for {

    // check for end of file
    line, err := reader.Read()
    if err == io.EOF {
        break
    } else if err != nil {
        log.Fatal(err)
    }

    // create label for entry
    track_label := createFilename(line[1], line[2])

    // see if already have this file
    own := haveTrack(track_label, have)

    // if don't have it, add to download queue list
    if !own {
      var entry track
      entry.url = line[0]
      entry.title = line[1]
      entry.artist = line[2]
      entry.album = line[3]
      entry.year = line[4]
      tracks = append(tracks, entry)
    }

  }

  return tracks
}

func haveTrack(test string, have []string) bool {
  for _, label := range have {
    if label == test {
      return true
    }
  }
  return false
}

func createFilename(title string, artist string) string {
  nospace := regexp.MustCompile(`\s`)
  nospec  := regexp.MustCompile(`[^\w\s]`)

  tmp := fmt.Sprintf("%s %s", artist, title)
  tmp = strings.ToLower(tmp)
  tmp = nospec.ReplaceAllString(tmp, "")
  tmp = nospace.ReplaceAllString(tmp, "_")

  return tmp + ".opus"
}
