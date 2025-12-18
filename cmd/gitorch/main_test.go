package main

import (
	"os/exec"
	"testing"
)

func CheckIfGitRepo() (bool, error) {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}

func TestCheckIfGitRepo(t *testing.T) {
	_, err := CheckIfGitRepo()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}
