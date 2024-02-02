package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/comatory/naive-nx/internal"
	"log"
)

var (
	stubborn  = flag.Bool("stubborn", false, "Force lint, type-check and test targets even if they are not present. Can result in failure.")
	baseRef   = flag.String("base-ref", "master", "Base ref to compare against. Defaults to master.")
	help      = flag.Bool("help", false, "Show help.")
	autoFetch = flag.Bool("autoFetch", true, "Automatically fetches from origin before running. Defaults to true.")
)

func showHelp() {
	fmt.Println("Usage: naive-nx [options]")
	fmt.Println("Options:")
	fmt.Println("  --stubborn: Force lint, type-check and test targets even if they are not present. Can result in failure.")
	fmt.Println("  --base-ref: Base ref to compare against. Defaults to master.")
	fmt.Println("  --autoFetch: Automatically fetches from origin before running. Defaults to true.")
	fmt.Println("  --help: Show help.")
	fmt.Println("Example: naive-nx (just runs with defaults)")
	fmt.Println("Example: git fetch some-other-origin && naive-nx --base-ref=some-other-origin/master (runs against different remote)")
}

func main() {
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *autoFetch {
		if err := internal.FetchOrigin(); err != nil {
			log.Fatal(errors.New("Could not fetch from origin. Is this project using git?"), err)
			return
		}
	}

	projectPathDescriptor, err := internal.GetNxProjectPaths()
	if err != nil {
		log.Fatal(errors.New("Could not find nx workspace. Is this an Nx project?"), err)
		return
	}

	touchedFiles, err := internal.GetFilepathsFromDiff(projectPathDescriptor.ProjectPath, *baseRef)
	if err != nil {
		log.Fatal(errors.New("Could not get filepaths from git diff. Is this project using git?"), err)
		return
	}

	projects := internal.GetAffectedNxProjectPaths(projectPathDescriptor.Paths, touchedFiles)

	if len(projects) == 0 {
		log.Println("No affected projects found. Is your master branch up-to-date?")
		return
	}

	projectDescriptors, err := internal.GetNxProjectDescriptors(projects)
	if err != nil {
		log.Fatal(errors.New("Could not get project names from paths. Is this an Nx project?"), err)
		return
	}

	log.Println("Preparation done. Running Nx commands...")

	err = internal.RunNxCommands(projectDescriptors, *stubborn)
	if err != nil {
		log.Fatal(errors.New("NX failed - check logs to see which target did not pass."), err)
		return
	}

	log.Println("OK!")
}
