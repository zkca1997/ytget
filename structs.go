package main

import (
	"log"
	"regexp"
	"strings"
	"path/filepath"
)

type track struct {
	id 					string
	public_url			string
	hidden_url			string
	contentLength		float64
	totalWrittenBytes	float64
	downloadLevel		float64
	title  				string
	artist 				string
	album  				string
	year   				string
	path				string
 	result				error
}

func (song *track) initTrack(music_dir string) {

	// if we are starting with a pathname
	if song.title == "" {

		// make sure necessary data is present
		if song.path == "" {
			log.Fatal("track init failed starting with path\n")
		}

		tmp := filepath.Base(song.path)
		song.title = strings.TrimRight(tmp, filepath.Ext(tmp))
		tmp = filepath.Dir(song.path)
		song.album = filepath.Base(tmp)
		tmp = filepath.Dir(tmp)
		song.artist = filepath.Base(tmp)
	}

	// if we are starting with metadata
	if song.path == "" {

		// make sure necessary data is present
		if (song.title == "" || song.artist == "" || song.album == "") {
			log.Fatal("track init failed starting with metadata\n")
		}

		// compile regex patterns
		nospace := regexp.MustCompile(`\s`)
		nospec := regexp.MustCompile(`[^\w\s]`)

		// normalize track title
		title := strings.ToLower(song.title)
		title = nospec.ReplaceAllString(title, "")
		title = nospace.ReplaceAllString(title, "_")

		// normalize track artist
		artist := strings.ToLower(song.artist)
		artist = nospec.ReplaceAllString(artist, "")
		artist = nospace.ReplaceAllString(artist, "_")

		// normalize track album
		album := strings.ToLower(song.album)
		album = nospec.ReplaceAllString(album, "")
		album = nospace.ReplaceAllString(album, "_")

		// return target location
		song.path = filepath.Join(music_dir, artist, album, title)
	}

	if song.id == "" {
		// make sure necessary data is present
		if (song.title == "" || song.artist == "") {
			log.Fatal("track failed to generate valid id")
		}

		// compile regex patterns
		nospace := regexp.MustCompile(`\s`)
		nospec := regexp.MustCompile(`[^\w\s]`)

		// normalize track id string
		id := song.artist + " " + song.title
		id = strings.ToLower(id)
		id  = nospec.ReplaceAllString(id, "")
		id  = nospace.ReplaceAllString(id, "_");

		song.id = id
	}
}