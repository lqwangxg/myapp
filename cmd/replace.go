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
		fmt.Printf("name=%s,  textfile=%s, directory=%s, suffix=%s \n", flags.name, flags.destFile, flags.destDir, flags.suffix)
		if flags.name == "" {
			fmt.Printf("--name=%s is required.", flags.name)
		}
		if flags.destFile == "" && flags.destDir == "" {
			fmt.Printf("--textfile=%s and --directory=%s can't be empty neither.", flags.destFile, flags.destDir)
		}
	},
}

func init() {
	regexCmd.AddCommand(replaceCmd)
}
