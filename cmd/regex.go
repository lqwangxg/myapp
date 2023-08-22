/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"
)

var wg sync.WaitGroup

// regexCmd represents the regex command
var regexCmd = &cobra.Command{
	Use:   "regex",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("regex %s called", flags.Action)
		if flags.Pattern != "" {
			if flags.Content != "" {
				Exec(NewRegexText(flags.Pattern, flags.Content))
			} else if flags.DestFile != "" {
				Exec(NewRegexFileByPattern(flags.Pattern, flags.RuleName, flags.DestFile))
			}
		} else if flags.RuleName != "" && flags.DestFile != "" {
			Exec(NewRegexFile(flags.RuleName, flags.DestFile))
		}
		if flags.RuleName != "" && flags.DestDir != "" {
			Exec(NewRegexDirectory(flags.RuleName, flags.DestDir))
		}
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(regexCmd)
	regexCmd.PersistentFlags().StringVarP(&flags.Action, "action", "a", "", "regex action match/replace")
	regexCmd.PersistentFlags().StringVarP(&flags.RuleName, "name", "n", "", "regex rule name which used to find yaml file or cache.")
	regexCmd.PersistentFlags().StringVarP(&flags.DestFile, "destFile", "f", "", "replace destination text file path")
	regexCmd.PersistentFlags().StringVarP(&flags.DestDir, "destDir", "d", "", "replace destination directory")
	regexCmd.PersistentFlags().StringVarP(&flags.Pattern, "pattern", "p", "", "regex pattern string")
	regexCmd.PersistentFlags().BoolVarP(&flags.ExportFlag, "export-flag", "", true, "export regex matches result flag by export templates.")

	regexCmd.PersistentFlags().StringVarP(&flags.IncludeSuffix, "include-suffix", "", "", "include pattern of dest file, default empty for all files")
	regexCmd.PersistentFlags().StringVarP(&flags.ExcludeSuffix, "exclude-suffix", "", "", "exclude pattern of dest file, default empty for all files")
}
