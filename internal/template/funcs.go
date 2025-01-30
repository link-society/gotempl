package template

import (
	"errors"
	"os"
	"path/filepath"
)

func ReadFile(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = nil
		}
		return false, err
	}

	return true, nil
}

func IsDir(path string) (bool, error) {
	s, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return s.IsDir(), nil
}

func Walk(path string) ([]string, error) {
	files, err := ReadDir(path)
	if err != nil {
		return nil, err
	}

	var result []string

	for _, file := range files {
		filePath := filepath.Join(path, file)
		result = append(result, file)

		if ok, _ := IsDir(filePath); ok {
			subFiles, err := Walk(filePath)
			if err != nil {
				return nil, err
			}

			for _, subFile := range subFiles {
				result = append(result, filepath.Join(file, subFile))
			}
		}
	}

	return result, nil
}

func ReadDir(path string) ([]string, error) {
	fileInfos, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var result = make([]string, len(fileInfos))

	for i, fileInfo := range fileInfos {
		result[i] = fileInfo.Name()
	}

	return result, nil
}

var funcs = map[string]interface{}{
	"isDir":        IsDir,
	"osIsDir":      IsDir,
	"readDir":      ReadDir,
	"osReadDir":    ReadDir,
	"readFile":     ReadFile,
	"osReadFile":   ReadFile,
	"walkDir":      Walk,
	"osWalkDir":    Walk,
	"fileExists":   Exists,
	"osFileExists": Exists,
}
