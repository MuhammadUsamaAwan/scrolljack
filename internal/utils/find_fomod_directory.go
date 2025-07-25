package utils

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func FindFomodDirectory(rootDir string) (fomodDir, moduleConfigPath string, err error) {
	err = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && strings.EqualFold(d.Name(), "fomod") {
			configPath := filepath.Join(path, "ModuleConfig.xml")
			if _, err := os.Stat(configPath); err == nil {
				fomodDir = path
				moduleConfigPath = configPath
				return filepath.SkipAll
			}
		}
		return nil
	})

	if err != nil {
		return "", "", err
	}

	return fomodDir, moduleConfigPath, nil
}
