package router

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Package struct {
	Path  string
	Files []string
}

func ScanPackages(dir string) ([]*Package, error) {
	var packages []*Package

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		name := info.Name()
		if strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" {
			return filepath.SkipDir
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			fmt.Printf("读取目录失败 %s: %v\n", path, err)
			return nil
		}

		var goFiles []string
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			fileName := entry.Name()
			if !strings.HasSuffix(fileName, ".go") {
				continue
			}

			if strings.HasSuffix(fileName, "_test.go") {
				continue
			}

			goFiles = append(goFiles, filepath.Join(path, fileName))
		}

		if len(goFiles) > 0 {
			packages = append(packages, &Package{
				Path:  path,
				Files: goFiles,
			})
		}

		return nil
	})

	return packages, err
}
