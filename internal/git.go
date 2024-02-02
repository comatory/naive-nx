package internal

import (
	"os"
	"os/exec"
	"path"
	"strings"
)

func ResolveRef(ref string) string {
	if ref == "" {
		return "master"
	}

	return ref
}

func FetchOrigin() error {
	cmd := exec.Command(
		"git",
		"fetch",
		"-v",
		"origin",
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()
}

func GetFilepathsFromDiff(projectPath string, ref string) ([]string, error) {
	branchRef := ResolveRef(ref)

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
