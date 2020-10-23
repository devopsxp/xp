/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"github.com/devopsxp/xp/pipeline"
	"github.com/spf13/cobra"
)

var (
	module string // 模块
	shell  string // shell 命令
	src    string // copy模块 src
	dest   string // copy模块 dest
)

// cliCmd represents the cli command
var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "命令行工具",
	Long: `指定模块进行单项功能使用. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: 完成数据Message.Data模型拼装
		data := map[string]interface{}{
			"module": module,
			"shell":  shell,
		}

		config := pipeline.DefaultPipeConfig("cli").
			WithInputName("cli").SetArgs(data).
			WithFilterName("shell").
			WithOutputName("console")

		p := pipeline.Of(*config)
		p.Init()
		p.Start()
		p.Exec()
		p.Stop()
	},
}

func init() {
	rootCmd.AddCommand(cliCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cliCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cliCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cliCmd.Flags().StringVarP(&module, "module", "m", "", "指定模块")
	cliCmd.Flags().StringVarP(&shell, "shell", "a", "", "执行命令")
}
