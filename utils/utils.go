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
	var exec string

	if shared {
		appsPath = "/usr/local/share/applications/"
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
	fmt.Println(_name, " ", savePath)
	savePath = filepath.Join(savePath, _name)
	fmt.Println(savePath)
	_path := filepath.Dir(path)
	fmt.Println(_path)
	executableName := filepath.Base(path)
	if filepath.Ext(path) == ".exe" {
		exec = "wine " + path
	} else {
		exec = path
	}

	var data string = "[Desktop Entry]\nVersion=1.1\n"
	data += "Type=Application\n"
	data += "Name=" + name + "\n"
	iconCheck, _k_ := ExtractEXEIcon(path)
	fmt.Println(_k_)
	if iconCheck {
		_icon := strings.TrimSuffix(executableName, filepath.Ext(executableName)) + ".png"
		data += "Icon=" + filepath.Join(_path, _icon) + "\n"
	}
	data += "Exec=" + exec + "\n"
	data += "Path=" + _path + "\n"
	data += "Actions=" + "\n"
	data += "Categories=Game;" + "\n"
	data += "Comment=" + desc + "\n"
	data += "Terminal=false\nStartupNotify=false\n"

	errr := os.WriteFile(savePath, []byte(data), 0644)
	if errr != nil {
		return fmt.Errorf("failed to save desktop entry: %w", errr)
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

func ExtractEXEIcon(exePath string) (bool, error) {
	outPath := filepath.Dir(exePath)
	err := os.MkdirAll(outPath, 0755) // Using os.MkdirAll to create the directory
	if err != nil {
		return false, fmt.Errorf("failed to create output directory: %v", err)
	}

	exe := filepath.Base(exePath)
	exeName := strings.TrimSuffix(exe, filepath.Ext(exe))
	cmd := exec.Command("wrestool", "-x", "-t", "14", exePath)

	outFile, err := os.Create(filepath.Join(outPath, exeName+".ico"))
	if err != nil {
		return false, fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return false, fmt.Errorf("failed to extract icon: %v", err)
	}

	cmd = exec.Command("convert", filepath.Join(outPath, exeName+".ico"), "-thumbnail", "32x32", filepath.Join(outPath, exeName+".png"))

	err = cmd.Run()
	if err != nil {
		return false, fmt.Errorf("failed to extract icon: %v", err)
	}

	return true, nil
}
