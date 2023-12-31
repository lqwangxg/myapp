package cmd

import (
	"io"
	"os"
)

func IsExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// isDirectory determines if a file represented
// by `path` is a directory or not
func IsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func ReadAll(filePath string) (string, error) {

	f, err := os.Open(filePath)
	check(err)
	defer f.Close()

	content, err := io.ReadAll(f) // 全部読み込んでくれる
	check(err)

	return string(content), nil
}
func WriteAll(filePath, content string) {
	f, err := os.Create(filePath)
	check(err)
	defer f.Close()

	_, err = f.WriteString(content)
	check(err)
	f.Sync()
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}
