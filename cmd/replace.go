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
		rs := NewRegexFromCmd()
		rs.Action = ReplaceAction
		//flags.RuleName
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
	regexCmd.AddCommand(replaceCmd)
}
