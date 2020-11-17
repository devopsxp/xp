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
	"strings"

	"github.com/devopsxp/xp/pipeline"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "抓取文件到管理机上",
	Long: `官方文档：https://docs.ansible.com/ansible/latest/modules/fetch_module.html#fetch-module
	eg: ./xp cli fetch 127.0.0.1 -u lxp -S /tmp/123 -D /tmp/333`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("Cli args: %v", args)
		if cliSrc == "" || cliDest == "" {
			log.Error("src or dest is not config")
			os.Exit(1)
		}

		log.Debugln("items", cliItem)
		var items []interface{}

		if cliItem != "" {
			items = []interface{}{}
			for _, x := range strings.Split(cliItem, ",") {
				items = append(items, x)
			}
		}

		data := map[string]interface{}{
			"host":        args,
			"remote_user": cliUser,
			"remote_pwd":  cliPwd,
			"remote_port": cliPort,
			"roles":       []interface{}{"fetch"},
			"stage":       []interface{}{"fetch"},
			"vars":        map[string]interface{}{},
			"hooks":       []interface{}{map[interface{}]interface{}{"type": "none"}},
			"config": []interface{}{map[interface{}]interface{}{
				"stage":      "fetch",
				"name":       "下载文件模块",
				"with_items": items,
				"copy": map[interface{}]interface{}{
					"src":  cliSrc,
					"dest": cliDest,
				},
			}},
		}

		config := pipeline.DefaultPipeConfig("fetch").
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
	cliCmd.AddCommand(fetchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	fetchCmd.Flags().StringVarP(&cliSrc, "src", "S", "", "远程目标主机文件 [不是目录]，批量eg: {{.item}}")
	fetchCmd.Flags().StringVarP(&cliDest, "dest", "D", "", "本地保存文件路径 [不是目录],批量eg: /tmp/{{.item}}")
	fetchCmd.Flags().StringVarP(&cliItem, "items", "I", "", "批量文件上传,eg: /tmp/1,/usr/kubectl./bin/docker")
}
