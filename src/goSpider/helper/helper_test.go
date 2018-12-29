package helper

import (
	"fmt"
	"testing"
)

func TestEnv(t *testing.T) {
	fmt.Println(Env())
}

func TestDownloadPath(t *testing.T) {
	fmt.Println(WorkspacePath())
	fmt.Println(WorkspacePath())
	fmt.Println(WorkspacePath())
}
