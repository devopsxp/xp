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
	"github.com/spf13/viper"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "基于模板的远程文件传输",
	Long: `使用go template进行文件模板解析并上传到远程目标主机上。
模板变量通过读取配置文件模块[vars]设置（全部改为小写字母）。
eg: ./xp cli template 127.0.0.1 -u lxp -S template.service.j2  -D /tmp/docker.service
example:
#==============================DEMO=================================
{{if .Status}} <font color="info">[成功] </font>{{ .title}} {{- else}}<font color="warning">[失败] </font> {{ .title}}  {{- end}} 
> 服务名称：<font color="comment">{{.serviceName}}</font> 
流水线名称：<font color="comment">{{.pipelineName}}</font> 
流水线ID：<font color="comment">{{.pipelineId}}</font> 
流水线地址：<font color="comment">{{.pipelineUrl}}</font> 
构建分支：<font color="comment">{{.branch}}</font> 
阶段：<font color="comment">{{.stage}}</font> 
构建人：<font color="comment">{{.user}}</font> 
COMMITID：<font color="comment">{{.tags}}</font> 
COMMITINFO: <font color="comment">{{.info}}</font> 

{{range $e := .alerts}}
## {{if eq $e.status "firing"}} [<font color="warning">{{$e.labels.severity}}</font>] {{ $e.labels.node}} [监控详情]({{$e.generatorURL}})

>> 告警阈值: <font color="warning">{{$e.labels.threshold_value}}</font> 
当前数值: <font color="warning">{{$e.annotations.current_value}}</font> 
开始时间: <font color="comment">{{$e.startsAt}}</font> 
结束时间: <font color="comment">{{$e.endsAt}}</font> 

{{- else}}# [<font color="info">{{$e.labels.severity}}</font>] {{ $e.labels.node}} [监控详情]({{$e.generatorURL}})

>> 告警阈值: <font color="info">{{$e.labels.threshold_value}}</font> 
当前数值: <font color="info">{{$e.annotations.current_value}}</font> 
开始时间: <font color="comment">{{$e.startsAt}}</font> 
结束时间: <font color="comment">{{$e.endsAt}}</font> 

{{- end}}

{{end}}
#==============================DEMO=================================`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("Cli args: %v Vars: %v", args, viper.GetStringMap("vars"))
		if cliSrc == "" || cliDest == "" {
			log.Error("src or dest is not config")
			os.Exit(1)
		}

		data := map[string]interface{}{
			"host":        args,
			"remote_user": cliUser,
			"remote_pwd":  cliPwd,
			"remote_port": cliPort,
			"roles":       []interface{}{"template"},
			"stage":       []interface{}{"template"},
			"vars":        viper.GetStringMap("vars"),
			"hooks":       []interface{}{map[interface{}]interface{}{"type": cliLogout}},
			"config": []interface{}{map[interface{}]interface{}{
				"stage": "template",
				"name":  "模板文件上传",
				"template": map[interface{}]interface{}{
					"src":  cliSrc,
					"dest": cliDest,
				},
			}},
		}

		config := pipeline.DefaultPipeConfig("template").
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
	cliCmd.AddCommand(templateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// templateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// templateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	templateCmd.Flags().StringVarP(&cliSrc, "src", "S", "", "本机模板文件")
	templateCmd.Flags().StringVarP(&cliDest, "dest", "D", "", "目标机上传路径")
}
