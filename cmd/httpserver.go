/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func startHttpServer() {
	// 创建一个多路复用器，用来根据不同的请求路径和方法调用不同的处理函数
	mux := http.NewServeMux()
	// 注册handleGet函数，用来处理"/"路径下的GET请求
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGet(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	// 注册handlePost函数，用来处理"/message"路径下的POST请求
	mux.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlePost(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// 创建一个服务器对象，设置监听地址和多路复用器
	server := &http.Server{
		Addr:    appContext.AppConfig.WebServer.Port,
		Handler: mux,
	}
	// 启动服务器，并监听端口8080上的请求
	fmt.Println("Server is running on port 8080")
	log.Fatal(server.ListenAndServe())
}

// 定义一个结构体，用来存储请求和响应的数据
type Message struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

// 定义一个处理GET请求的函数，返回一个欢迎信息
func handleGet(w http.ResponseWriter, r *http.Request) {
	// 设置响应头的内容类型为JSON
	w.Header().Set("Content-Type", "application/json")
	// 创建一个Message对象，赋值给变量msg
	msg := Message{
		Name: "Bing",
		Text: "Welcome to my restapi server!",
	}
	// 将msg对象转换为JSON格式的字节切片，赋值给变量data
	data, err := json.Marshal(msg)
	// 如果转换出错，打印错误信息，并返回状态码500
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// 将data字节切片写入响应体中
	w.Write(data)
}

// 定义一个处理POST请求的函数，返回一个回复信息
func handlePost(w http.ResponseWriter, r *http.Request) {
	// 设置响应头的内容类型为JSON
	w.Header().Set("Content-Type", "application/json")
	// 创建一个Message对象，赋值给变量msg
	msg := Message{}
	// 从请求体中读取JSON格式的数据，并解析到msg对象中
	err := json.NewDecoder(r.Body).Decode(&msg)
	// 如果解析出错，打印错误信息，并返回状态码400
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	// 关闭请求体
	defer r.Body.Close()
	// 打印接收到的数据
	fmt.Printf("Received: %s: %s\n", msg.Name, msg.Text)
	// 修改msg对象的Text属性，添加一个回复信息
	msg.Text = "Hello, " + msg.Name + ". Thank you for your message."
	// 将msg对象转换为JSON格式的字节切片，赋值给变量data
	data, err := json.Marshal(msg)
	// 如果转换出错，打印错误信息，并返回状态码500
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// 将data字节切片写入响应体中
	w.Write(data)
}
