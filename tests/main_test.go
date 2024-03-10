package tests

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestSmt(t *testing.T) {
	b, err := exec.Command("ls", "-la").Output()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(b))
}
