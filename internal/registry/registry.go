package registry

import (
	"golang.org/x/sys/windows/registry"
)

func Read(k registry.Key, path string, name string) (string, error) {

	k, err := registry.OpenKey(k, path, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}

	s, _, err := k.GetStringValue(name)
	if err != nil {
		_ = k.Close()
		return "", err
	}

	_ = k.Close()
	return s, nil
}

func Write(k registry.Key, path string, name string, value string) error {

	k, err := registry.OpenKey(k, path, registry.WRITE)
	if err != nil {
		return err
	}

	err = k.SetStringValue(name, value)
	if err != nil {
		_ = k.Close()
		return err
	}

	_ = k.Close()
	return nil
}
