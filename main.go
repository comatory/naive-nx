package main

import (
	"errors"
	"github.com/comatory/naive-nx/internal"
	"log"
)

func main() {
	projectPathDescriptor, err := internal.GetNxProjectPaths()
	if err != nil {
		log.Fatal(errors.New("Could not find nx workspace. Is this an Nx project?"), err)
		return
	}

	touchedFiles, err := internal.GetFilepathsFromDiff(projectPathDescriptor.ProjectPath)
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

	err = internal.RunNxCommands(projectDescriptors)
	if err != nil {
		log.Fatal(errors.New("NX failed - check logs to see which target did not pass."), err)
		return
	}

	log.Println("OK!")
}
