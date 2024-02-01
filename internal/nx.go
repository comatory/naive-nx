package internal

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type WorkspaceFileDescriptor struct {
	filePath    string
	projectPath string
}

type ProjectPathDescriptor struct {
	ProjectPath string
	Paths       []string
}

type WorkspaceJson struct {
	Schema   string `json:"$schema"`
	Projects map[string]string
}

type ProjectJson struct {
	Name string
  Targets map[string]interface{}
}

type ProjectDescriptor struct {
  Name string
  Targets []string
}

const MAX_DEPTH = 10

func parseWorkspaceJson(file *os.File) (*WorkspaceJson, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	jsonData := &WorkspaceJson{}
	if err := json.Unmarshal(data, jsonData); err != nil {
		return nil, err
	}

	return jsonData, nil
}

func getPathToNxProjects(workspaceFilePath string) (map[string]string, error) {
	file, err := os.Open(workspaceFilePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	jsonData, jsonErr := parseWorkspaceJson(file)

	if jsonErr != nil {
		return nil, jsonErr
	}

	return jsonData.Projects, nil
}

func validateWorkspaceFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, jsonErr := parseWorkspaceJson(file)

	if jsonErr != nil {
		return jsonErr
	}

	if jsonData.Schema != "" {
		if strings.Contains(jsonData.Schema, "nx/schemas/workspace-schema.json") {
			return nil
		} else {
			return errors.New("$schema key does not contain the expected value.")
		}
	} else {
		return errors.New("$schema key not found.")
	}
}

func findNxRoot() (*WorkspaceFileDescriptor, error) {
	currentDir, err := os.Getwd()

	if err != nil {
		return &WorkspaceFileDescriptor{}, err
	}

	counter := 0

	for {
		if counter > MAX_DEPTH {
			return &WorkspaceFileDescriptor{}, errors.New("Max depth reached")
		}

		filePath := filepath.Join(currentDir, "workspace.json")

		_, err := os.Stat(filePath)
		if err == nil {
			workspaceValidationErr := validateWorkspaceFile(filePath)

			if workspaceValidationErr != nil {
				return &WorkspaceFileDescriptor{}, workspaceValidationErr
			}

			return &WorkspaceFileDescriptor{
				filePath:    filePath,
				projectPath: currentDir,
			}, nil // File found
		}

		if currentDir == filepath.Dir(currentDir) {
			return &WorkspaceFileDescriptor{}, os.ErrNotExist
		}

		currentDir = filepath.Dir(currentDir)
		counter++
	}
}

func GetNxProjectPaths() (*ProjectPathDescriptor, error) {
	workspaceFileDescriptor, validationErr := findNxRoot()

	if validationErr != nil {
		return nil, validationErr
	}

	projects, err := getPathToNxProjects(workspaceFileDescriptor.filePath)

	if err != nil {
		return nil, err
	}

	var paths []string
	for _, value := range projects {
		paths = append(paths, workspaceFileDescriptor.projectPath+string(filepath.Separator)+value)
	}

	return &ProjectPathDescriptor{
		Paths:       paths,
		ProjectPath: workspaceFileDescriptor.projectPath,
	}, nil
}

func GetAffectedNxProjectPaths(projectPaths []string, filePaths []string) []string {
	var affectedPaths []string

	for _, projectPath := range projectPaths {
		for _, filePath := range filePaths {
			if strings.HasPrefix(filePath, projectPath) {
				affectedPaths = append(affectedPaths, projectPath)
			}
		}
	}

	var uniquePaths []string

	for _, affectedPath := range affectedPaths {
		found := false
		for _, uniquePath := range uniquePaths {
			if affectedPath == uniquePath {
				found = true
				break
			}
		}

		if !found {
			uniquePaths = append(uniquePaths, affectedPath)
		}
	}

	return uniquePaths
}

func extractValidTargets(targets map[string]interface{}) []string {
  var validTargets []string

  for key := range targets {
      if key == "lint" || key == "type-check" || key == "test" {
        validTargets = append(validTargets, key)
      }
  }

  return validTargets
}

func GetNxProjectDescriptors(projectRootPaths []string) ([]*ProjectDescriptor, error) {
	var descriptors []*ProjectDescriptor

	for _, projectRootPath := range projectRootPaths {
		projectJsonPath := projectRootPath + string(filepath.Separator) + "project.json"

		file, err := os.Open(projectJsonPath)
		if err != nil {
			return nil, err
		}

		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		projectJson := &ProjectJson{}
		if err := json.Unmarshal(data, projectJson); err != nil {
			return nil, err
		}

		descriptors = append(descriptors, &ProjectDescriptor{
        Name: projectJson.Name,
        Targets: extractValidTargets(projectJson.Targets),
    })
	}

	return descriptors, nil
}
