/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var name, textfile, directory, suffix string

// replaceCmd represents the replaceTF command
var replaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "replace text file or text files under directory by pattern name",
	Long: ` replace text file or text files under directory by pattern name. 
	      pattern name connects to a json or configMap which includes rules of pattern/replacement/skipRules.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("regex replace called")
		fmt.Printf("name=%s,  textfile=%s, directory=%s, suffix=%s \n", name, textfile, directory, suffix)
		if name == "" {
			fmt.Printf("--name=%s is required.", name)
		}
		if textfile == "" && directory == "" {
			fmt.Printf("--textfile=%s and --directory=%s can't be empty neither.", textfile, directory)
		}
		// if(textfile !=""){
		// }
	},
}

// func replaceFile(name string, textfile string) {

// }

// func replaceDirectory(name string, directory string) {

// }

func init() {
	regexCmd.AddCommand(replaceCmd)

	replaceCmd.Flags().StringVarP(&name, "name", "n", "", "regex replace pattern name")
	replaceCmd.Flags().StringVarP(&textfile, "destFile", "f", "", "replace destination text file path")
	replaceCmd.Flags().StringVarP(&directory, "destDir", "d", "", "replace destination directory")
	replaceCmd.Flags().StringVarP(&suffix, "suffix", "s", "", "replace destination file suffix, default empty")
}
