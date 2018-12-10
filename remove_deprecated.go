package main

import (
	"os"
	"fmt"
	"runtime"
	"strings"
)

func remove_deprecated(del []*track, music_dir string) {

	// print all tracks marked for removal
	fmt.Println("\nTracks Marked for Removal from Library:")
	fmt.Println("---------------------------------------")
	for _, entry := range del {
		fmt.Printf("%s.ogg\n", entry.id)
	}

	// prompt user input and hang until valid response
	for {
		fmt.Print("\nConfirm Removal of Listed Tracks [y/N]: ")
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
				for _, song := range del {
					err := os.Remove(song.path)
					if err != nil { fmt.Println(err) }
				}
				return;
			case "N":
				return;
			default:
		}
	}
}
