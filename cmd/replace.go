/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// replaceCmd represents the replaceTF command
var replaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "replace text file or text files under directory by pattern name",
	Long: ` replace text file or text files under directory by pattern name. 
	      pattern name connects to a json or configMap which includes rules of pattern/replacement/skipRules.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("regex replace called")

		if flags.Pattern != "" && flags.Content != "" {
			handler := NewRegexText(flags.Pattern, flags.Content)
			reger.Execute(handler)
		}
		if flags.RuleName != "" && flags.DestFile != "" {
			handler := NewRegexFile(flags.RuleName, flags.DestFile)
			reger.Execute(handler)
		}
		if flags.RuleName != "" && flags.DestDir != "" {
			handler := NewRegexDirectory(flags.RuleName, flags.DestDir)
			reger.Execute(handler)
		}
	},
}

func init() {
	regexCmd.AddCommand(replaceCmd)
}
