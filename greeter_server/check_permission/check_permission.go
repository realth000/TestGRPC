package check_permission

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	permitDirs []string
)

func GetPermittedDir() []string {
	return permitDirs
}

func LoadPermission(loadPath string) error {
	if loadPath != "" {
		confFile, err := os.Open(loadPath)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(confFile)
		for scanner.Scan() {
			// TODO: Check if path valid here.
			path, err := filepath.Abs(scanner.Text())
			if err != nil {
				continue
			}
			permitDirs = append(permitDirs, path)
		}
	}

	// PermitFiles download files in current directory if no permit paths loaded.
	if len(permitDirs) == 0 {
		currentPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return errors.New("no permission can load")
		}
		permitDirs = append(permitDirs, currentPath)
	}
	return nil
}

func CheckPathPermission(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	for _, s := range permitDirs {
		if strings.HasPrefix(absPath, s) && filepath.Base(s) == filepath.Base(absPath) {
			return true
		}
	}
	return false
}
