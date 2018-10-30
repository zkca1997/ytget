package main

import (
	"os"
	"fmt"
	"runtime"
	"strings"
	"path/filepath"
)

func remove_deprecated(del []*track, music_dir string) error {

	// print all tracks marked for removal
	fmt.Println("\nTracks Marked for Removal from Library:")
	fmt.Println("---------------------------------------")
	for _, entry := range del {
		fmt.Println(entry)
	}

	// prompt user input and hang until valid response
	for {
		fmt.Println("\nConfirm Removal of Listed Tracks [y/N]: ")
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
				deleteFiles(del, music_dir)
				break;
			case "N":
				os.Exit(0)
			default:
		}

	}
}

func deleteFiles(list []*track, music_dir string) {
	for _, song := range list {
		err := os.Remove(filepath.Join(music_dir, song.path))
		if err != nil {
			fmt.Printf("failed to delete: %s\n", song.path)
		}
	}
}