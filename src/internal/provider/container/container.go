package container

import (
	"fmt"
	"os/exec"
)

func ExecContainer() error {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to execute 'go version': %w", err)
	}

	fmt.Println(string(output))

	return nil
}
