package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"path/filepath"
)

type track struct {
	url    string
	title  string
	artist string
	album  string
	year   string
}

func main() {

	music_dir := filepath.Join(os.Getenv("HOME"), "Music")
	have := getCurrentTracks(music_dir)
	want, data := readManifest(filepath.Join(music_dir, ".meta", "manifest.csv"))

	get, del := findDiff(have, want)

	// print list of tracks to download
	if len(get) > 0 {
		fmt.Printf("\nAdding to Queue:\n----------------\n")
		for _, entry := range get {
			fmt.Printf("%s by %s\n", data[entry].title, data[entry].artist)
		}
	}

	// print list of tracks to delete
	if len(del) > 0 {
		fmt.Printf("\nDeleting from Library:\n----------------------\n")
		for _, entry := range del {
			fmt.Println(entry)
		}
	}

	// prompt for user consent
	fmt.Println("\nPress Enter to confirm changes (Ctrl-C to Cancel): ")
	var input string
	fmt.Scanln(&input)

	queue_file := music_dir + "/.meta/queue.csv"
	writeQueue(get, data, queue_file)
	deleteFiles(del, music_dir)
}

func deleteFiles(tracks []string, directory string) {
	for _, track := range tracks {
		err := os.Remove(directory + "/" + track)
		if err != nil {
			fmt.Printf("failed to delete: %s\n", track)
		}
	}
}

func findDiff(have []string, want []string) ([]string, []string) {
	/* this is very inefficient, awaiting sorting optimization */

	var get []string
	for _, a := range want {
		match := false
		for _, b := range have {
			if strings.Compare(a, b) == 0 {
				match = true
				break
			}
		}
		if match == false {
			get = append(get, a)
		}
	}

	var del []string
	for _, a := range have {
		match := false
		for _, b := range want {
			if strings.Compare(a, b) == 0 {
				match = true
				break
			}
		}
		if match == false {
			del = append(del, a)
		}
	}

	return get, del
}

func writeQueue(tracks []string, data map[string]track, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, entry := range tracks {
		var data = []string{
			data[entry].url,
			data[entry].title,
			data[entry].artist,
			data[entry].album,
			data[entry].year}
		err := writer.Write(data)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getCurrentTracks(music_path string) []string {

	pattern := filepath.Join(music_path, "*", "*", "*")
	longtracks, err := filepath.Glob(pattern)
	if err != nil { log.Fatalf("Fatal Error: %s\n", err) }
	var tracks []string

	for _, member := range longtracks {
		title := filepath.Base(member)
		tmp := filepath.Dir(member)
		album := filepath.Base(tmp)
		tmp = filepath.Dir(tmp)
		artist := filepath.Base(tmp)
		foo := filepath.Join(artist, album, title)
		tracks = append(tracks, foo)
	}

	return tracks
}

func readManifest(filename string) ([]string, map[string]track) {

	// compile regex patterns
	nospace := regexp.MustCompile(`\s`)
	nospec := regexp.MustCompile(`[^\w\s]`)

	// open file
	manifestFile, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error opening %s: %s\n", filename, err)
	}

	// open CSV parser
	reader := csv.NewReader(bufio.NewReader(manifestFile))
	data := make(map[string]track)
	var filelist []string

	for {
		// check for end of file
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		// filename added to list
		title := strings.ToLower(line[1])
		title = nospec.ReplaceAllString(title, "")
		title = nospace.ReplaceAllString(title, "_")

		artist := strings.ToLower(line[2])
		artist = nospec.ReplaceAllString(artist, "")
		artist = nospace.ReplaceAllString(artist, "_")

		album := strings.ToLower(line[3])
		album = nospec.ReplaceAllString(album, "")
		album = nospace.ReplaceAllString(album, "_")

		file := filepath.Join(artist, album, title + ".opus")
		filelist = append(filelist, file)

		// read the metadata and append to slice
		var entry track
		entry.url = line[0]
		entry.title = line[1]
		entry.artist = line[2]
		entry.album = line[3]
		entry.year = line[4]
		data[file] = entry
	}

	return filelist, data
}
