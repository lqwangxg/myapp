/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	//"context"
	"github.com/spf13/cobra"
	//"github.com/redis/go-redis/v9"
)

var ruleName, ruleStr string

// ruleCmd represents the rule command
var ruleCmd = &cobra.Command{
	Use:   "rule",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rule called")
		fmt.Printf("origin flag: ruleName=%s,  ruleStr=%s \n", ruleName, ruleStr)
		//ReadYaml()
		//setValue(ruleName, ruleStr + " ===>.suffix")
		//newRule := getValue(ruleName)
		//fmt.Printf("after setValue ruleName=%s,  old-ruleStr=%s, new-ruleStr \n", ruleName, ruleStr,newRule)
	},
}

// var client = redis.NewClient(&redis.Options{
// 	Addr:	  "localhost:6379",
// 	Password: "", // no password set
// 	DB:		  0,  // use default DB
// });
// func getValue(key string) string {
// 	ctx := context.Background()
// 	val, err := client.Get(ctx, key).Result()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("getValue: key=%s,  rule=%s \n", key, val)
// 	return val
// }
// func setValue(key string, value string) {
// 	ctx := context.Background()
// 	err := client.Set(ctx, key, value, 0).Err()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("setValue: key=%s,  rule=%s \n", key, value)
// }

func init() {
	regexCmd.AddCommand(ruleCmd)
	ruleCmd.Flags().StringVarP(&ruleName, "ruleName", "", "", "regex replace pattern name")
	ruleCmd.Flags().StringVarP(&ruleStr, "ruleStr", "", "", "regex replace rule")
}
