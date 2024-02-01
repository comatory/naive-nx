package internal

import (
	"testing"
)

func TestGetAffectedNxProjectPaths(t *testing.T) {
	affectedPaths := GetAffectedNxProjectPaths(
		[]string{"libs/lib-a", "libs/lib-b", "libs/lib-c"},
		[]string{"libs/lib-a/src/lib-a.ts", "libs/lib-b/src/lib-b.ts"})

	if len(affectedPaths) != 2 {
		t.Errorf("Expected 2 affected paths, got %d", len(affectedPaths))
	}

	if affectedPaths[0] != "libs/lib-a" {
		t.Errorf("Expected libs/lib-a, got %s", affectedPaths[0])
	}

	if affectedPaths[1] != "libs/lib-b" {
		t.Errorf("Expected libs/lib-b, got %s", affectedPaths[1])
	}
}

func TestGetAffectedNxProjectPathsNoMatch(t *testing.T) {
	affectedPaths := GetAffectedNxProjectPaths(
		[]string{"libs/lib-a", "libs/lib-b", "libs/lib-c"},
		[]string{"libs/lib-d/src/lib-d.ts", "libs/lib-e/src/lib-e.ts"})

	if len(affectedPaths) != 0 {
		t.Errorf("Expected 0 affected paths, got %d", len(affectedPaths))
	}
}

func TestGetAffectedNxProjectPathsNoPaths(t *testing.T) {
	affectedPaths := GetAffectedNxProjectPaths(
		[]string{},
		[]string{"libs/lib-d/src/lib-d.ts", "libs/lib-e/src/lib-e.ts"})

	if len(affectedPaths) != 0 {
		t.Errorf("Expected 0 affected paths, got %d", len(affectedPaths))
	}
}

func TestGetAffectedNxProjectPathsNoFiles(t *testing.T) {
	affectedPaths := GetAffectedNxProjectPaths(
		[]string{"libs/lib-a", "libs/lib-b", "libs/lib-c"},
		[]string{})

	if len(affectedPaths) != 0 {
		t.Errorf("Expected 0 affected paths, got %d", len(affectedPaths))
	}
}

func TestGetAffectedNxProjectPathsNoPathsNoFiles(t *testing.T) {
	affectedPaths := GetAffectedNxProjectPaths(
		[]string{},
		[]string{})

	if len(affectedPaths) != 0 {
		t.Errorf("Expected 0 affected paths, got %d", len(affectedPaths))
	}
}

func TestGetAffectedNxProjectUniquePaths(t *testing.T) {
	affectedPaths := GetAffectedNxProjectPaths(
		[]string{"libs/lib-a", "libs/lib-b", "libs/lib-c"},
		[]string{"libs/lib-a/src/lib-a.ts", "libs/lib-a/src/lib-a.ts"})

	if len(affectedPaths) != 1 {
		t.Errorf("Expected 1 affected paths, got %d", len(affectedPaths))
	}

	if affectedPaths[0] != "libs/lib-a" {
		t.Errorf("Expected libs/lib-a, got %s", affectedPaths[0])
	}
}
