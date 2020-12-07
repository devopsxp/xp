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

var (
	imageIE      string // 镜像文件
	imageArgs    string // 参数
	imageCommand string // 命令
)

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "基于docker容器进行命令执行",
	Long:  `本机docker执行pipeline命令`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("Cli args: %v", args)
		if imageIE == "" {
			log.Error("未检测到执行命令,请确认！ [eg: ./xp cli image -i java:8 -a \"-v /tmp:/data -u root \" java -version && for i in {1..10};do echo `date`;sleep 1;done]")
			os.Exit(1)
		}

		tmpc := []interface{}{}
		for _, c := range args {
			tmpc = append(tmpc, c)
		}

		data := map[string]interface{}{
			"host":        []string{"127.0.0.1"},
			"remote_user": cliUser,
			"remote_pwd":  cliPwd,
			"remote_port": cliPort,
			"roles":       []interface{}{"image"}, // shell role and stage
			"vars":        map[string]interface{}{},
			"hooks":       []interface{}{map[interface{}]interface{}{"type": cliLogout}},
			"stage":       []interface{}{"image"},
			"config": []interface{}{map[interface{}]interface{}{
				"stage":   "image",
				"name":    "Image模块",
				"image":   imageIE,
				"args":    []interface{}{imageArgs},
				"command": tmpc,
			}},
		}

		config := pipeline.DefaultPipeConfig("image").
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
	cliCmd.AddCommand(imageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	imageCmd.Flags().StringVarP(&imageIE, "image", "i", "", "镜像名")
	imageCmd.Flags().StringVarP(&imageArgs, "args", "a", "", "docker启动参数，如：-v /tmp:/data -e USER=xp -u root")
	imageCmd.Flags().StringVarP(&imageCommand, "command", "c", "", "docker执行命令，可以理解为CMD")
}
