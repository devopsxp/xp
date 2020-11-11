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

var (
	cliSrc  string // 源地址
	cliDest string // 目标路径
	cliItem string // 多文件传输
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "远程传输文件",
	Long: `将src文件传输到远程目标主机dest路径，
eg: ./xp cli copy 127.0.0.1 -u lxp -S /tmp/123 -D /tmp/333`,
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
			"roles":       []interface{}{"copy"},
			"stage":       []interface{}{"copy"},
			"vars":        map[string]interface{}{},
			"hooks":       []interface{}{map[interface{}]interface{}{"type": "none"}},
			"config": []interface{}{map[interface{}]interface{}{
				"stage":      "copy",
				"name":       "上传文件模块",
				"with_items": items,
				"copy": map[interface{}]interface{}{
					"src":  cliSrc,
					"dest": cliDest,
				},
			}},
		}

		config := pipeline.DefaultPipeConfig("copy").
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
	cliCmd.AddCommand(copyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// copyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// copyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	copyCmd.Flags().StringVarP(&cliUser, "user", "u", "root", "远程主机执行用户，默认：root")
	copyCmd.Flags().StringVarP(&cliSrc, "src", "S", "", "原路径，批量eg: {{.item}}")
	copyCmd.Flags().StringVarP(&cliDest, "dest", "D", "", "目的路径,批量eg: /tmp/{{.item}}")
	copyCmd.Flags().StringVarP(&cliItem, "items", "I", "", "批量文件上传,eg: /tmp/1,/usr/kubectl./bin/docker")
}
