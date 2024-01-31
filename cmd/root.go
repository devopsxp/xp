/*
Copyright © 2020 NAME lflxp 382023823@qq.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	_ "github.com/devopsxp/xp/module"
	"github.com/spf13/viper"
)

var cfgFile string
var debug bool
var islog bool
var cliUser string   // 远程执行用户
var cliPwd string    // 远程执行用户密码
var cliPort int      // 远程ssh 端口
var cliLogout string // 日志输出格式

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "devopsxp",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/devopsxp.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "是否打印debug日志")
	rootCmd.PersistentFlags().BoolVarP(&islog, "log", "l", false, "是否文件输出")
	rootCmd.PersistentFlags().StringVarP(&cliUser, "user", "u", "root", "远程主机执行用户，默认：root")
	rootCmd.PersistentFlags().StringVarP(&cliPwd, "pwd", "p", "", "远程主机用户密码，默认：")
	rootCmd.PersistentFlags().IntVarP(&cliPort, "port", "P", 22, "远程主机ssh端口，默认：22")
	rootCmd.PersistentFlags().StringVarP(&cliLogout, "logout", "L", "count", "日志格式：console|none|email|wechat|count|k8shook")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// 获取项目的执行路径
		home, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		// Search config in home directory with name ".devopsxp" (without extension).
		viper.AddConfigPath(home)       // 设置读取文件的路径
		viper.SetConfigName("devopsxp") // 设置读取的文件名
		viper.SetConfigType("yaml")     // 设置文件的类型
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		slog.Debug("Using config file:", viper.ConfigFileUsed())
	} else {
		slog.Error("Using config file error", "ERROR", err.Error())
	}
}
