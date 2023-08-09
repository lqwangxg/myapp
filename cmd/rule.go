/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	//"context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	//"github.com/redis/go-redis/v9"
)

//var ruleName, ruleStr string

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
		//fmt.Printf("origin flag: ruleName=%s,  ruleStr=%s \n", ruleName, ruleStr)
		ReadYaml()
		//setValue(ruleName, ruleStr + " ===>.suffix")
		//newRule := getValue(ruleName)
		//fmt.Printf("after setValue ruleName=%s,  old-ruleStr=%s, new-ruleStr \n", ruleName, ruleStr,newRule)
	},
}

func ReadYaml() {
	viper.SetConfigType("yaml")
	viper.SetConfigType("rule/html-input.yaml")
	//viper.SetConfigFile("config/rules/html-input.yaml")
	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())
	viper.ReadInConfig()

	// if viper.IsSet("config.global_params") {
	// 	global_params := viper.Get("config.global_params") //.([]string)
	// 	fmt.Println("config.global_params:", global_params)
	// } else {
	// 	fmt.Println(" config.global_params not set.")
	// }
	// if viper.IsSet("config.spec_chars") {
	// 	//params := global_params.([]map[string]string)
	// 	spec_chars := viper.Get("config.spec_chars") //.([]map[string]string)
	// 	fmt.Println("config.spec_chars:", spec_chars)
	// } else {
	// 	fmt.Println(" config.spec_chars not set.")
	// }
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
}
