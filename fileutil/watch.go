package fileutil

import (
	"os"
	"time"
)

// Watch watches a file for changes and calls the action if there are changes.
func Watch(path string, action func(*os.File) error) error {
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	lastMod := stat.ModTime()

	ticker := time.Tick(500 * time.Millisecond)
	for {
		select {
		case <-ticker:
			stat, err = os.Stat(path)
			if err != nil {
				return err
			}
			if stat.ModTime().After(lastMod) {
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				if err := action(file); err != nil {
					return err
				}
				lastMod = stat.ModTime()
			}
		}
	}
}
