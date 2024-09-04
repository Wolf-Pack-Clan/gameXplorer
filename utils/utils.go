package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type GameEntry struct {
	Name    string
	Desc    string
	Exec    string
	ExecDir string
	Icon    string
}

func ListGames() ([]GameEntry, error) {
	var games []GameEntry

	// Paths to check for .desktop files
	paths := []string{
		"/usr/share/applications/",
		"~/.local/share/applications/",
	}

	for _, path := range paths {
		// Expand ~ to home directory
		if strings.HasPrefix(path, "~/") {
			path = filepath.Join(os.Getenv("HOME"), path[2:])
		}

		err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(filePath, ".desktop") {
				file, err := os.Open(filePath)
				if err != nil {
					return err
				}
				defer file.Close()

				var name, _desc, exec, _execdir, _icon string
				var isGame bool

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := scanner.Text()
					if strings.HasPrefix(line, "Name=") {
						name = strings.TrimPrefix(line, "Name=")
					} else if strings.HasPrefix(line, "Comment=") {
						_desc = strings.TrimPrefix(line, "Comment=")
					} else if strings.HasPrefix(line, "Exec=") {
						exec = strings.TrimPrefix(line, "Exec=")
					} else if strings.HasPrefix(line, "Path=") {
						_execdir = strings.TrimPrefix(line, "Path=")
					} else if strings.HasPrefix(line, "Icon=") {
						_icon = strings.TrimPrefix(line, "Icon=")
					} else if strings.HasPrefix(line, "Categories=") && strings.Contains(line, "Game;") {
						isGame = true
					}
				}

				if isGame && name != "" && exec != "" {
					games = append(games, GameEntry{Name: name, Desc: _desc, Exec: exec, ExecDir: _execdir, Icon: _icon})
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return games, nil
}

func SaveGame(name string, desc string, path string, shared bool) error {
	var appsPath string
	var savePath string

	if shared {
		appsPath = "/usr/share/applications"
	} else {
		appsPath = "~/.local/share/applications"
	}

	if strings.HasPrefix(appsPath, "~/") {
		savePath = filepath.Join(os.Getenv("HOME"), appsPath[2:])
	} else {
		savePath = appsPath
	}

	_name := strings.ToLower(strings.Replace(name, " ", "-", -1))
	_name += ".desktop"
	err := os.WriteFile(savePath, []byte(_name), 0644)
	if err != nil {
		return fmt.Errorf("failed to save path to file: %w", err)
	}
	fmt.Println("savePath:", savePath)
	return nil
}

func IsWineInstalled() bool {
	cmd := exec.Command("which", "wine")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func IsLutrisInstalled() bool {
	cmd := exec.Command("which", "lutris")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func ExtractEXEIcon(exePath string) error {
	outPath := filepath.Dir(exePath)
	err := exec.Command("mkdir", "-p", outPath).Run()
	if err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	cmd := exec.Command("wrestool", "-x", "--output="+outPath, "-t14", exePath)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to extract icon: %v", err)
	}
	return nil
}
