package main

import (
	"os"
	"fmt"
	"flag"
	"strings"
	"runtime"
	"path/filepath"
)

var (
	music_dir	string
	meta_file	string
)

func init() {

	default_music_dir  := filepath.Join(os.Getenv("HOME"), "Music")
	default_meta_file  := filepath.Join(default_music_dir, ".meta", "manifest.csv")

	flag.StringVar(&music_dir, "d", default_music_dir, "root of the music directory")
	flag.StringVar(&meta_file, "m", default_meta_file, "manifest file of music library")
}

func main() {

	// parse cli flags
	flag.Parse()

	// scan the meta file and library for updates
	get, del := scan(music_dir, meta_file)

	// if deprecated files, prompt for removal
	if len(del) > 0 { remove_deprecated(del, music_dir) }

	// download new songs
	// print all tracks marked for download
	fmt.Println("\nTracks Marked for Download:")
	fmt.Println("-----------------------------")
	for _, entry := range get {
		fmt.Printf("%s.ogg\n", entry.id)
	}

	// prompt user input and hang until valid response
	for {
		fmt.Print("\nConfirm Download of Listed Tracks [y/N]: ")
		var input string
		fmt.Scanln(&input)

		// trim raw STDIN input for native OS
		if runtime.GOOS == "windows" {
			input = strings.TrimRight(input, "\r\n")
		} else {
			input = strings.TrimRight(input, "\n")
		}

		switch input {
			case "y":
				sync(get)
				return;
			case "N":
				return;
			default:
		}

	}
}

func sync(Targets []*track) {

	// create buffered channels
	jobs  := make(chan *track, len(Targets))
	pipe1 := make(chan *track, len(Targets))
	pipe2 := make(chan *track, len(Targets))
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
		fmt.Printf("Job Error:\t%s\n", result)
	}
}
