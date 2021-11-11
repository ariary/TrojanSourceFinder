package excludelist

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetExcludelist(excludelistFilename string) (*map[string]struct{}, error) {
	excludelist := make(map[string]struct{})
	if excludelistFilename == "" {
		return &excludelist, nil
	}

	excludelistFile, err := os.Stat(excludelistFilename)
	if err != nil {
		return nil, err
	}

	if excludelistFile.IsDir() {
		return nil, fmt.Errorf("excludefile %s must not be a directory", excludelistFilename)
	}

	readfile, err := os.Open(excludelistFilename)
	if err != nil {
		return nil, err
	}
	defer readfile.Close()

	scanner := bufio.NewScanner(readfile)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		excludelist[filepath.Clean(line)] = struct{}{}
	}

	return &excludelist, scanner.Err()
}
