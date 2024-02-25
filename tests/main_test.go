package tests

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestSmt(t *testing.T) {
	s, err := exec.LookPath("Skillbox - DevOps-инженер. Основы")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s)
	b, err := exec.Command("ls", "-la").Output()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(b))
}
