/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
			rs := NewRegex(pattern)
			rs.MatchText(flags.Content)
		}
		if flags.DestFile != "" {
			MatchFile(pattern, flags.DestFile)
		}
		if flags.DestDir != "" {
			MatchDiretory(pattern, flags.DestDir)
		}
	},
}

func MatchFile(pattern, filePath string) {
	// exit if file not exists
	if !IsExists(filePath) {
		return
	}
	if flags.IncludeSuffix != "" {
		re := NewRegex(flags.IncludeSuffix)
		if !re.IsMatch(filePath) {
			return
		}
	}
	if flags.ExcludeSuffix != "" {
		re := NewRegex(flags.ExcludeSuffix)
		if re.IsMatch(filePath) {
			return
		}
	}
	if buffer, err := ReadAll(filePath); err == nil {
		rs := NewRegex(pattern)
		rs.Result.Params["filePath"] = filePath
		rs.MatchText(buffer)
	}
}

func MatchDiretory(pattern, dirPath string) {
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
		fullPath := filepath.Join(dirPath, file.Name())
		ok, err := IsDir(fullPath)
		if err == nil {
			if ok {
				MatchDiretory(pattern, fullPath)
			} else {
				MatchFile(pattern, fullPath)
			}
		}
	}
}

func (rs *Regex) MatchText(content string) {
	rs.ScanMatches(content)
	rs.Close()

	if value, ok := rs.Result.Params["filePath"]; ok {
		log.Printf("filePath: %s", value)
	}
	//export matches contents
	if export := rs.ExportMatches(flags.Template); export != "" {
		log.Printf("export: %s", export)
	} else {
		log.Printf("export: empty. no matches")
	}

}

// func ReplaceText(pattern, content, matchTemplate string) {
// 	rs := NewRegex(pattern)
// 	rs.ScanMatches(content)
// 	//rs.ExportMatches(matchTemplate)
// 	rs.Close()

// 	if flags.Template != "" {
// 		log.Print(rs.ToString())
// 	}
// 	//rs.SplitMatch()
// 	//log.Print(rs.ToString())
// 	//rs.log()
// }

func init() {
	regexCmd.AddCommand(matchCmd)
	TestReplaceMap()
}
