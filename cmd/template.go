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

	"github.com/spf13/cobra"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "基于模板的远程文件传输",
	Long: `使用go template进行文件模板解析并上传到远程目标主机上。example:
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
		fmt.Println("template called")
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
}
