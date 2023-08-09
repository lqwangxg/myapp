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

	regexCmd.PersistentFlags().StringVarP(&flags.destFile, "destFile", "f", "", "replace destination text file path")
	regexCmd.PersistentFlags().StringVarP(&flags.destDir, "destDir", "d", "", "replace destination directory")
	regexCmd.PersistentFlags().StringVarP(&flags.pattern, "pattern", "p", "", "regex pattern string")
	regexCmd.PersistentFlags().StringVarP(&flags.origin, "content", "c", "", "input content string")

	regexCmd.PersistentFlags().StringVarP(&flags.name, "name", "n", "", "regex replace pattern name")
	regexCmd.PersistentFlags().StringVarP(&flags.suffix, "suffix", "s", "", "replace destination file suffix, default empty")
}
