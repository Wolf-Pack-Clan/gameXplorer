package utils

import (
	"bufio"
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

func IsWineInstalled() bool {
	cmd := exec.Command("which", "wine")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
