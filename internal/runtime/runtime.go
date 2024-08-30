package runtime

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

const PREFIX string = "."

type Runtime struct {
	Name     string
	Path     string
	Configs  []string
	Versions []string
}

func Get(wd string) ([]Runtime, error) {
	globs, err := filepath.Glob(filepath.Join(wd, PREFIX+"*"))
	if err != nil {
		return nil, err
	}

	runtimes := make([]Runtime, 0, len(globs))

	for _, path := range globs {
		name, _ := strings.CutPrefix(filepath.Base(path), PREFIX)

		configs, err := readConfigs(path)

		if err != nil {
			continue
		}

		versions, err := readVersions(path)

		if err != nil {
			continue
		}

		runtimes = append(runtimes, Runtime{
			Name:     name,
			Path:     path,
			Configs:  configs,
			Versions: versions,
		})
	}

	return runtimes, err
}

func readConfigs(path string) ([]string, error) {
	file := filepath.Join(path, "configs.json")

	read, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var configs []string
	err = json.Unmarshal(read, &configs)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

func readVersions(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	versions := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		versions = append(versions, entry.Name())
	}

	return versions, nil
}
