/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "myapp",
	Short: "read params and flags in cui",
}

var config AppConfig
var appContext AppContext = *NewContext(&config)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&flags.ConfigFile, "config", ".myapp.yaml", "config file (default is .myapp.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	//fmt.Println("root initConfig called ", cfgFile)
	if flags.ConfigFile != "" {
		// Use config file from the flag.
		//viper.SetConfigFile(flags.ConfigFile)
		appContext.LoadFile(flags.ConfigFile)
	}
	// // Find home directory.
	// home, err := os.UserHomeDir()
	// cobra.CheckErr(err)

	// // Search config in home directory with name ".myapp" (without extension).
	// viper.AddConfigPath(home)
	// viper.SetConfigType("yaml")
	// viper.AutomaticEnv() // read in environment variables that match

	appContext.LoadDirectory(appContext.AppConfig.RuleDir)
}
