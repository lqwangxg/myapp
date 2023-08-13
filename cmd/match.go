/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

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
		rs := NewRegex(pattern)
		rs.Action = MatchAction
		if flags.Content != "" {
			rs.MatchText(flags.Content)
		}
		if flags.DestFile != "" {
			rs.ProcFile(flags.DestFile)
		}
		if flags.DestDir != "" {
			rs.ProcDir(flags.DestDir)
		}
	},
}

func init() {
	regexCmd.AddCommand(matchCmd)
}
