package cmd

import (
	"io"
	"os"
)

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
