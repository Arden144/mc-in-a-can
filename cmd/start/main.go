package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/Arden144/paperupdate/pkg/update"
)

var version string
var memory string
var jar string

//go:embed config/*
var embeds embed.FS

var args = []string{
	"-XX:+UseG1GC",
	"-XX:+ParallelRefProcEnabled",
	"-XX:MaxGCPauseMillis=200",
	"-XX:+UnlockExperimentalVMOptions",
	"-XX:+DisableExplicitGC",
	"-XX:+AlwaysPreTouch",
	"-XX:G1NewSizePercent=30",
	"-XX:G1MaxNewSizePercent=40",
	"-XX:G1HeapRegionSize=8M",
	"-XX:G1ReservePercent=20",
	"-XX:G1HeapWastePercent=5",
	"-XX:G1MixedGCCountTarget=4",
	"-XX:InitiatingHeapOccupancyPercent=15",
	"-XX:G1MixedGCLiveThresholdPercent=90",
	"-XX:G1RSetUpdatingPauseTimePercent=5",
	"-XX:SurvivorRatio=32",
	"-XX:+PerfDisableSharedMem",
	"-XX:MaxTenuringThreshold=1",
	"-Dusing.aikars.flags=https://mcflags.emc.gs",
	"-Daikars.new.flags=true",
}

func init() {
	flag.StringVar(&version, "version", "latest", "")
	flag.StringVar(&memory, "memory", "4G", "")
	flag.StringVar(&jar, "jar", "paper.jar", "")
	flag.Parse()
}

func exists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

func main() {
	var err error
	dir, err := embeds.ReadDir("config")
	if err != nil {
		fmt.Printf("failed to read virtual filesystem: %v", err)
		return
	}
	for _, original := range dir {
		path := original.Name()

		if exists, err := exists(path); err == nil && exists {
			continue
		} else if err != nil {
			fmt.Printf("failed to check for file %v: %v", path, err)
			return
		}

		data, err := embeds.ReadFile("config/" + path)
		if err != nil {
			fmt.Printf("failed to read virtual file %v: %v", "config/"+path, err)
			return
		}

		file, err := os.Create(path)
		if err != nil {
			fmt.Printf("failed to create file %v: %v", path, err)
			return
		}

		if _, err = file.Write(data); err != nil {
			file.Close()
			os.Remove(path)
			fmt.Printf("failed to write file %v: %v", path, err)
			return
		}

		if err = file.Close(); err != nil {
			os.Remove(path)
			fmt.Printf("failed to close file %v: %v", path, err)
			return
		}
	}

	fmt.Println("Checking for updates")

	if build, err := update.TryUpdate(version, jar); err == nil {
		fmt.Printf("Updated to %v build %v\n", version, build)
	} else if errors.Is(err, update.ErrLatest) {
		fmt.Println("Already up to date")
	} else {
		fmt.Printf("Failed to update: %v", err)
	}

	args = append(args, "-Xms"+memory, "-Xmx"+memory, "-jar", jar, "nogui")
	cmd := exec.Command("java", args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to start minecraft: %v", err)
	}
}
