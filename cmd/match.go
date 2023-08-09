/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// matchCmd represents the match command
var matchCmd = &cobra.Command{
	Use:   "match",
	Short: "match string by regrep",
	Long:  ` match string by regrep pattern, and replace string if parameter --replace is set.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("regex match called")
		pattern := flags.pattern
		if pattern == "" {
			return
		}
		if flags.origin != "" {
			MatchText(pattern, flags.origin)
		}
		if flags.destFile != "" {
			MatchFile(pattern, flags.destFile, flags.suffix)
		}
		if flags.destDir != "" {
			MatchDiretory(pattern, flags.destFile, flags.suffix)
		}
	},
}

func MatchFile(pattern, filePath, suffix string) {
	// exit if file not exists
	if !IsExists(filePath) {
		return
	}

	if buffer, err := ReadAll(filePath); err == nil {
		MatchText(pattern, buffer)
	}
}

func MatchDiretory(pattern, dirPath, suffix string) {
	if !IsExists(dirPath) {
		log.Printf("dirPath is not found. dirPath=%s", dirPath)
		return
	}
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, file := range files {
		ok, err := IsDir(file.Name())
		if err == nil {
			if ok {
				MatchDiretory(pattern, file.Name(), suffix)
			} else {
				MatchFile(pattern, file.Name(), suffix)
			}
		}
	}
}
func MatchText(pattern, content string) {
	content = beforeMatch(content)
	matches := Matches(pattern, content)
	matches.restore()
	matches.log()
}

// func ReplaceText(pattern string, content string, replacement string) {

// 	match, _ := regexp.MatchString(pattern, content)
// 	fmt.Println(match)

// 	// r, _ := regexp.Compile(pattern)

// 	// fmt.Println(r.MatchString(content))

// 	// fmt.Println(r.FindString(content))

// 	// fmt.Println(r.FindStringIndex(content))

// 	// fmt.Println(r.FindStringSubmatch(content))

// 	// fmt.Println(r.FindStringSubmatchIndex(content))

// 	// fmt.Println(r.FindAllString(content, -1))

// 	// fmt.Println(r.FindAllStringSubmatchIndex(content, -1))

// 	// fmt.Println(r.FindAllString(content, 2))

// 	// fmt.Println(r.Match([]byte(content)))

// 	r := regexp.MustCompile(pattern)
// 	fmt.Println(r)

// 	fmt.Println(r.ReplaceAllString(content, replacement))

// 	in := []byte(content)
// 	out := r.ReplaceAllFunc(in, bytes.ToUpper)
// 	fmt.Println(string(out))
// }

func init() {
	regexCmd.AddCommand(matchCmd)
}
