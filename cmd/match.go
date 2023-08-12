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
		pattern := flags.Pattern
		if pattern == "" {
			return
		}
		if flags.Content != "" {
			MatchText(pattern, flags.Content)
		}
		if flags.DestFile != "" {
			MatchFile(pattern, flags.DestFile, flags.Suffix)
		}
		if flags.DestDir != "" {
			MatchDiretory(pattern, flags.DestFile, flags.Suffix)
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
	rs := NewRegex(pattern)
	rs.ScanMatches(content)
	rs.Close()

	if flags.Template != "" {
		log.Print(rs.ToString())
	}
	//rs.SplitMatch()
	//log.Print(rs.ToString())
	//rs.log()
}

func init() {
	regexCmd.AddCommand(matchCmd)
	TestReplaceMap()
}
