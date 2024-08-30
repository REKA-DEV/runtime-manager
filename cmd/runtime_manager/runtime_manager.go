package main

import (
	"errors"
	registry2 "github.com/REKA-DEV/runtime-manager/internal/registry"
	"github.com/REKA-DEV/runtime-manager/internal/runtime"
	"github.com/REKA-DEV/runtime-manager/internal/selector"
	"golang.org/x/sys/windows/registry"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}

	runtimes, err := runtime.Get(wd)
	if err != nil {
		log.Fatal(err)
		return
	}

	/*for _, runtime := range runtimes {
		println("name: " + runtime.Name)
		println("path: " + runtime.Path)
		println("configs: ")
		for _, config := range runtime.Configs {
			println("         ", config)
		}
		println("versions: ")
		for _, version := range runtime.Versions {
			println("         ", version)
		}
		println("========")
	}*/

	path, err := registry2.Read(registry.CURRENT_USER, "Environment", "Path")
	if err != nil {
		log.Fatal(err)
		return
	}

	path = strings.Trim(path, ";")
	isChanged := false

	for _, r := range runtimes {
		for _, config := range r.Configs {
			p := filepath.Join(wd, r.Name, config)
			if strings.Contains(path, p) {
				continue
			}
			path = path + ";" + p
			isChanged = true
		}
	}

	if isChanged {
		err2 := registry2.Write(registry.CURRENT_USER, "Environment", "Path", path)
		if err2 != nil {
			log.Fatal(err2)
			return
		}
	}

	runtimeMenu := selector.New[runtime.Runtime]("choose runtime")

	for _, r := range runtimes {
		runtimeMenu.Add(r.Name, &r)
	}

	runtimeChoice, err := runtimeMenu.Run()
	if err != nil {
		log.Fatal(err)
		return
	}

	versionMenu := selector.New[string]("choose version")

	for _, v := range runtimeChoice.Versions {
		versionMenu.Add(v, &v)
	}

	versionChoice, err := versionMenu.Run()
	if err != nil {
		log.Fatal(err)
		return
	}

	symlink := filepath.Join(wd, runtimeChoice.Name)

	_, err = os.Stat(symlink)
	if !errors.Is(err, os.ErrNotExist) {
		err = os.Remove(symlink)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	err = os.Symlink(filepath.Join(runtimeChoice.Path, *versionChoice), filepath.Join(wd, runtimeChoice.Name))
	if err != nil {
		log.Fatal(err)
		return
	}
}
