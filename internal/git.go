package internal

import (
	"os/exec"
	"path"
	"strings"
)

func GetFilepathsFromDiff(projectPath string, ref string) ([]string, error) {
	branchRef := ref

	if ref == "" {
		branchRef = "master"
	}

	cmd := exec.Command(
		"git",
		"diff",
		"--name-only",
		branchRef)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, err
	}

	paths := strings.Split(string(output), "\n")
	absolutePaths := make([]string, len(paths))

	for i, p := range paths {
		absolutePaths[i] = path.Join(projectPath, p)
	}

	return absolutePaths, nil
}
