/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// regexCmd represents the regex command
var regexCmd = &cobra.Command{
	Use:   "regex",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("regex called")
	},
}

func init() {
	rootCmd.AddCommand(regexCmd)
	regexCmd.PersistentFlags().StringVarP(&flags.RuleName, "name", "n", "", "regex rule name which used to find yaml file or cache.")
	regexCmd.PersistentFlags().StringVarP(&flags.DestFile, "destFile", "f", "", "replace destination text file path")
	regexCmd.PersistentFlags().StringVarP(&flags.DestDir, "destDir", "d", "", "replace destination directory")
	regexCmd.PersistentFlags().StringVarP(&flags.Pattern, "pattern", "p", "", "regex pattern string")

	regexCmd.PersistentFlags().StringVarP(&flags.IncludeSuffix, "include-suffix", "", "", "include pattern of dest file, default empty for all files")
	regexCmd.PersistentFlags().StringVarP(&flags.ExcludeSuffix, "exclude-suffix", "", "", "exclude pattern of dest file, default empty for all files")
	regexCmd.PersistentFlags().StringVarP(&flags.ReplaceTemplate, "replace-template", "", "", "replace template for replacing input string")
	regexCmd.PersistentFlags().StringVarP(&flags.ExportTemplate, "export-template", "", "", "export template for exporting matches string")
}
