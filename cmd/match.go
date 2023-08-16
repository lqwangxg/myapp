/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

// import (
// 	"fmt"

// 	"github.com/spf13/cobra"
// )

// // matchCmd represents the match command
// var matchCmd = &cobra.Command{
// 	Use:   "match",
// 	Short: "match string by regrep",
// 	Long:  ` match string by regrep pattern, and replace string if parameter --replace is set.`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("regex match called")
// 		if flags.Pattern != "" && flags.Content != "" {
// 			handler := NewRegexText(flags.Pattern, flags.Content)
// 			reger.Execute(handler)
// 		}
// 		if flags.RuleName != "" && flags.DestFile != "" {
// 			handler := NewRegexFile(flags.RuleName, flags.DestFile)
// 			reger.Execute(handler)
// 		}
// 		if flags.RuleName != "" && flags.DestDir != "" {
// 			handler := NewRegexDirectory(flags.RuleName, flags.DestDir)
// 			reger.Execute(handler)
// 		}
// 	},
// }

// func init() {
// 	regexCmd.AddCommand(matchCmd)
// }

// func ProcFile(filePath string) {
// 	// exit if file not exists
// 	if !IsExists(filePath) {
// 		return
// 	}
// 	if flags.IncludeSuffix != "" {
// 		if !IsMatchString(flags.IncludeSuffix, filePath) {
// 			return
// 		}
// 	}
// 	if rs.Rule.IncludeFile != "" {
// 		if !IsMatchString(rs.Rule.IncludeFile, filePath) {
// 			return
// 		}
// 	}
// 	if flags.ExcludeSuffix != "" {
// 		if IsMatchString(flags.ExcludeSuffix, filePath) {
// 			return
// 		}
// 	}
// 	if rs.Rule.ExcludeFile != "" {
// 		if IsMatchString(rs.Rule.ExcludeFile, filePath) {
// 			return
// 		}
// 	}

// 	if buffer, err := ReadAll(filePath); err == nil {
// 		//rs.FromFile = filePath
// 		//log.Printf("Matching file: %s", rs.FromFile)
// 		//ProcString(buffer)
// 	}
// }
