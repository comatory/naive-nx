package internal

import (
	"os/exec"
	"strings"
	"path"
)

func GetFilepathsFromDiff(projectPath string) ([]string, error) {
	cmd := exec.Command(
		"git",
		"diff",
		"--name-only",
		"master")

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
