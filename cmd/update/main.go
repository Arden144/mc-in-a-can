package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Arden144/paperupdate/pkg/hash"
	"github.com/Arden144/paperupdate/pkg/paper"
	"github.com/Arden144/paperupdate/pkg/req"
)

func main() {
	project := paper.Project("paper")

	versions, err := project.Versions()
	if err != nil {
		fmt.Print(err)
		return
	}

	qs := []*survey.Question{
		{
			Name:   "File",
			Prompt: &survey.Input{Message: "Enter the file name"},
			Validate: func(val interface{}) error {
				if str, ok := val.(string); !ok || !strings.HasSuffix(str, ".jar") {
					return fmt.Errorf("the file name must end with .jar")
				}
				return nil
			},
		},
		{
			Name: "Version",
			Prompt: &survey.Select{
				Message: "Select a version",
				Options: versions,
			},
		},
	}

	options := struct{ File, Version string }{}
	err = survey.Ask(qs, &options)
	if err != nil {
		fmt.Print(err)
		return
	}

	download, err := project.DownloadInfo(options.Version)
	if err != nil {
		fmt.Print(err)
		return
	}

	if hash, err := hash.FromFile(options.File); err == nil {
		if hash == download.Sha256 {
			fmt.Print("No update is required")
			return
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		fmt.Print(err)
		return
	}

	written, err := req.GetFile(download.Url, options.File)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("Checksum %v\n", download.Sha256)
	fmt.Printf("Saved %v bytes", written)
}
