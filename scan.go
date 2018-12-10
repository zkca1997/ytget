package main

import (
	"os"
	"io"
	"log"
	"bufio"
	"encoding/csv"
	"path/filepath"

)

func scan(music_dir string, meta_file string) ([]*track, []*track) {
	
	have := scanLibrary(music_dir)	// scan library for current files
	want := readManifest(meta_file, music_dir)	// scan manifest for desired files

	// get list of songs we need to remove
	var del []*track
	for id, _ := range(have) {
		if _, ok := want[id]; !ok {
			tmp := have[id]
			del = append(del, &tmp)
		}
	}

	// get list of songs we need to download
	var get []*track
	for id, _ := range(want) {
		if _, ok := have[id]; !ok {
			tmp := want[id]
			get = append(get, &tmp)
		}
	}

	return get, del
}

func scanLibrary(music_path string) map[string]track {

	pattern := filepath.Join(music_path, "*", "*", "*.ogg")

	track_path, err := filepath.Glob(pattern)
	if err != nil { log.Fatalf("Fatal Error: %s\n", err) }

	library := make(map[string]track)

	for _, member := range track_path {

		song := track{ path: member }
		song.initTrack(music_path)

		library[song.id] = song
	}

	return library
}

func readManifest(filename string, music_path string) map[string]track {

	// open file
	manifestFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening %s: %s\n", filename, err)
	}

	list := make(map[string]track)

	// open CSV parser
	reader := csv.NewReader(bufio.NewReader(manifestFile))

	for {
		// check for end of file
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		// populate Youtube struct
		entry := track{
			public_url:	line[0],
			title:		line[1],
			artist:		line[2],
			album:		line[3],
			year:		line[4],
		}
		entry.initTrack(music_path)

		// map the entry into library
		list[entry.id] = entry
	}

	return list
}