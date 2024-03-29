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
	"fmt"
	"log/slog"
	"os"

	"github.com/devopsxp/xp/pipeline"
	"github.com/spf13/cobra"
)

var (
	cliShell     string // shell 命令
	cliTerminial bool   // 是否交互式执行命令
)

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "执行远程命令",
	Long:  `example: ./xp cli shell 127.0.0.1-20 192.168.50.1-10 -a "zsh" -T -L console`,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug(fmt.Sprintf("Cli args: %v", args))
		if cliShell == "" {
			slog.Error("未检测到执行命令,请确认！ [eg: ./xp cli shell 127.0.0.1-20 -a \"hostname\"]")
			os.Exit(1)
		} else if len(args) == 0 {
			slog.Error("未检测到目标主机，请确认！ [eg: ./xp cli shell 127.0.0.1-20 -a \"hostname\"]")
			os.Exit(1)
		}
		// TODO: 完成数据Message.Data模型拼装
		data := map[string]interface{}{
			"host":        args,
			"remote_user": cliUser,
			"remote_pwd":  cliPwd,
			"remote_port": cliPort,
			"roles":       []interface{}{"shell"}, // shell role and stage
			"vars":        map[string]interface{}{},
			"hooks":       []interface{}{map[interface{}]interface{}{"type": cliLogout}},
			"stage":       []interface{}{"shell"},
			"config":      []interface{}{map[interface{}]interface{}{"stage": "shell", "name": "Shell命令模块", "shell": cliShell, "terminial": cliTerminial}},
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
	cliCmd.AddCommand(shellCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// shellCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// shellCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	shellCmd.Flags().StringVarP(&cliShell, "shell", "a", "", "执行命令")
	shellCmd.Flags().BoolVarP(&cliTerminial, "terminial", "T", false, "是否执行交互式操作")
}
