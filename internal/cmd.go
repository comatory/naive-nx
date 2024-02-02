package internal

import (
	"os"
	"os/exec"
	"strings"
)

var defaultTargets = []string{"lint", "type-check", "test"}

func RunNxCommands(projectDescriptors []*ProjectDescriptor, stubbornMode bool) error {
	errors := make([]error, 0)

	if stubbornMode {
		projectNames := make([]string, len(projectDescriptors))

		for i, projectDescriptor := range projectDescriptors {
			projectNames[i] = projectDescriptor.Name
		}

		for _, target := range defaultTargets {
			cmd := exec.Command("yarn", "nx", "run-many", "--target", target, "--projects", strings.Join(projectNames, ","))

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Env = os.Environ()

			err := cmd.Run()
			if err != nil {
				errors = append(errors, err)
			}
		}
	} else {
		for _, projectDescriptor := range projectDescriptors {
			for _, target := range projectDescriptor.Targets {
				cmd := exec.Command("yarn", "nx", target, projectDescriptor.Name)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Env = os.Environ()

				err := cmd.Run()
				if err != nil {
					errors = append(errors, err)
				}
			}
		}
	}

	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}
