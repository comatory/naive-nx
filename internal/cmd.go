package internal

import (
	"os"
	"os/exec"
)

func RunNxCommands(projectDescriptors []*ProjectDescriptor) error  {
	errors := make([]error, 0)

	for _, projectDescriptor := range projectDescriptors {
		for _, target := range projectDescriptor.Targets {
			cmd := exec.Command("nx", target, projectDescriptor.Name)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Env = os.Environ()

			err := cmd.Run(); if err != nil {
				errors = append(errors, err)
			}
		}
	}

	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}
