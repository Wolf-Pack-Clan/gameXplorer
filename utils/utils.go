package utils

import (
	"os/exec"
)

// CheckIfWineInstalled checks if Wine is installed on the system
func IsWineInstalled() bool {
	cmd := exec.Command("which", "wine")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
