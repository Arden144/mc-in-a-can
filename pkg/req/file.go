package req

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
)

func GetFile(url, path string) (written int64, err error) {
	defer func() {
		if r := recover(); r != nil {
			written, err = 0, fmt.Errorf("download exception: %v", r)
		}
	}()

	// Request to GET the given url
	r, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("http request failed: %v", err)
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("http GET failed: %v", r.Status)
	}

	// Open a temp file for the download
	tmp, err := os.Create("paperupdate.tmp")
	if err != nil {
		return 0, fmt.Errorf("tempfile creation failed: %v", err)
	}
	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()

	// Open a progress bar
	progress := progressbar.DefaultBytes(r.ContentLength, path)
	defer progress.Close()

	// Write to the temp file and progress bar
	written, err = io.Copy(io.MultiWriter(tmp, progress), r.Body)
	if err != nil {
		return 0, fmt.Errorf("download failed: %v", err)
	}

	// Close the temp file
	if err = tmp.Close(); err != nil {
		return 0, fmt.Errorf("tempfile close failed: %v", err)
	}

	// Move the temp file to the given path
	if err = os.Rename(tmp.Name(), path); err != nil {
		return 0, fmt.Errorf("failed to save to %v: %v", path, err)
	}

	return
}
