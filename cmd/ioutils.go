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

// func main() {
// 	files, err := os.ReadDir(".")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for _, file := range files {
// 		log.Println(file.Name())
// 	}
// }

func ReadAll(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	content, err := io.ReadAll(f) // 全部読み込んでくれる
	if err != nil {
		return "", err
	}

	//fmt.Println(string(content))
	return string(content), nil
}
