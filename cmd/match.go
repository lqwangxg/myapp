/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"fmt"
	"regexp"
	"github.com/spf13/cobra"
)

// matchCmd represents the match command
var matchCmd = &cobra.Command{
	Use:   "match",
	Short: "match string by regrep",
	Long: ` match string by regrep pattern, and replace string if parameter --replace is set.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("regex match called")
		
		pattern, _ := cmd.Flags().GetString("pattern")
		content, _ := cmd.Flags().GetString("content")
		replace, _ := cmd.Flags().GetString("replace")
		
		fmt.Printf("args=%v, args.length=%d pattern=%s,  content=%s, replace=%s \n", args, len(args), pattern, content, replace)
	
		MatchText(pattern, content)
		ReplaceText(pattern, content, replace)
	},
}
func MatchText(pattern string, content string) {

    match, _ := regexp.MatchString(pattern, content)
    fmt.Println(match)

    r, _ := regexp.Compile(pattern)

    fmt.Println(r.MatchString(content))

    fmt.Println(r.FindString(content))

    fmt.Println(r.FindStringIndex(content))

    fmt.Println(r.FindStringSubmatch(content))

    fmt.Println(r.FindStringSubmatchIndex(content))

    fmt.Println(r.FindAllString(content, -1))

    fmt.Println(r.FindAllStringSubmatchIndex(content, -1))

    fmt.Println(r.FindAllString(content, 2))

    fmt.Println(r.Match([]byte(content)))

	//FindStringSubmatch
    //格納されるデータはスライスで入っているので、一つ一つ取り出すことも可能。
    //indexだけ取り出したいなど
    //スライスで取得される
    fss := r.FindStringSubmatch("/index/test")
    fmt.Println(fss, fss[0], fss[1], fss[2], len(fss))
    //>>[/index/test index test] /index/test index test 3
    //スライスで取り出せる


    //改行がある場合を検知する場合^ $ではなく、\A \zを使う
    //セキュリティ上望ましい。
    r3:= regexp.MustCompile(`\A/(index|detail|create)/([a-zA-Z0-9]+)\z`)
    fs3 := r3.FindString("/index/test")
    fmt.Println(fs3)//len(fs3)

}

func ReplaceText(pattern string, content string, replacement string) {

    match, _ := regexp.MatchString(pattern, content)
    fmt.Println(match)

    // r, _ := regexp.Compile(pattern)

    // fmt.Println(r.MatchString(content))

    // fmt.Println(r.FindString(content))

    // fmt.Println(r.FindStringIndex(content))

    // fmt.Println(r.FindStringSubmatch(content))

    // fmt.Println(r.FindStringSubmatchIndex(content))

    // fmt.Println(r.FindAllString(content, -1))

    // fmt.Println(r.FindAllStringSubmatchIndex(content, -1))

    // fmt.Println(r.FindAllString(content, 2))

    // fmt.Println(r.Match([]byte(content)))

    r := regexp.MustCompile(pattern)
    fmt.Println(r)

    fmt.Println(r.ReplaceAllString(content, replacement))

    in := []byte(content)
    out := r.ReplaceAllFunc(in, bytes.ToUpper)
    fmt.Println(string(out))
}

func init() {
	regexCmd.AddCommand(matchCmd)
	matchCmd.Flags().StringP("pattern", "p","", "regex pattern string")
	matchCmd.Flags().StringP("content", "c","", "input content string")
	matchCmd.Flags().StringP("replace", "r","", "regex replace string")
	
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// matchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// matchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
