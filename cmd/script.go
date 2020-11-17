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
	"os"

	"github.com/devopsxp/xp/pipeline"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var script string

// scriptCmd represents the script command
var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "远程执行脚本",
	Long:  `eg: ./xp cli script 127.0.0.1 -a "test.sh"`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("Cli args: %v", args)
		if script == "" {
			log.Error("未指定需要执行的脚本路径")
			os.Exit(1)
		} else if len(args) == 0 {
			log.Error("未检测到目标主机，请确认！ [eg: ./xp cli script 127.0.0.1-20 -a test.sh]")
			os.Exit(1)
		}

		data := map[string]interface{}{
			"host":        args,
			"remote_user": cliUser,
			"remote_pwd":  cliPwd,
			"remote_port": cliPort,
			"roles":       []interface{}{"script"},
			"stage":       []interface{}{"script"},
			"vars":        map[string]interface{}{},
			"hooks":       []interface{}{map[interface{}]interface{}{"type": "none"}},
			"config": []interface{}{map[interface{}]interface{}{
				"stage":    "script",
				"name":     "脚本服务",
				"script":   script,
				"terminal": true,
			}},
		}

		config := pipeline.DefaultPipeConfig("systemd").
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
	cliCmd.AddCommand(scriptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scriptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scriptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	scriptCmd.Flags().StringVarP(&script, "script", "a", "", "脚本路径")
}
