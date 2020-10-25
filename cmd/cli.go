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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	module string // 模块
	shell  string // shell 命令
	src    string // copy模块 src
	dest   string // copy模块 dest
	user   string // 远程用户
)

// cliCmd represents the cli command
var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "命令行工具",
	Long: `指定模块进行单项功能使用. For example:

copy 远程文件传输
template 模板文件传输
shell 远程shell命令执行`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("Cli args: %v", args)
		// TODO: 完成数据Message.Data模型拼装
		data := map[string]interface{}{
			"host":        args,
			"remote_user": user,
			"roles":       []interface{}{"shell"}, // shell role and stage
			"vars":        map[string]interface{}{},
			"hooks":       []interface{}{map[interface{}]interface{}{"type": "console"}},
			"stage":       []interface{}{"shell"},
			"config":      []interface{}{map[interface{}]interface{}{"stage": "shell", "name": "Running", module: shell, "src": src, "dest": dest}},
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
	cliCmd.Flags().StringVarP(&src, "src", "S", "", "copy source路径")
	cliCmd.Flags().StringVarP(&dest, "dest", "D", "", "dest 目标路径")
	cliCmd.Flags().StringVarP(&user, "user", "u", "root", "远程主机执行用户，默认：root")
}
